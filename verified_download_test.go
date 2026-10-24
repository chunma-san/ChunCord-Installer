/*
 * SPDX-License-Identifier: GPL-3.0
 * Based on Vencord Installer and Equilotl.
 * ChunCord integration maintained by chunma-san.
 */
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fixtureTransport func(*http.Request) (*http.Response, error)

func (f fixtureTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestVerifiedDownloadAndReplace(t *testing.T) {
	data := "verified client fixture"
	sum := sha256.Sum256([]byte(data))
	checksum := hex.EncodeToString(sum[:])
	release := &GithubRelease{}
	asset := struct {
		Name        string `json:"name"`
		Size        int64  `json:"size"`
		DownloadURL string `json:"browser_download_url"`
	}{}
	prefix := "https://github.com/chunma-san/ChunCord/releases/download/test/"
	for _, name := range []string{"desktop.asar", "SHA256SUMS.txt"} {
		item := asset
		item.Name = name
		item.Size = int64(len(data))
		item.DownloadURL = prefix + name
		release.Assets = append(release.Assets, item)
	}
	corrupt := true
	oldClient := downloadClient
	downloadClient = &http.Client{Transport: fixtureTransport(func(request *http.Request) (*http.Response, error) {
		body := data
		if strings.HasSuffix(request.URL.Path, "SHA256SUMS.txt") {
			body = checksum + "  desktop.asar\n"
		} else if corrupt {
			body = "broken data"
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})}
	t.Cleanup(func() { downloadClient = oldClient })
	directory := t.TempDir()
	destination := filepath.Join(directory, "chuncord.asar")
	if err := os.WriteFile(destination, []byte("previous"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := downloadVerifiedArtifact(release, "desktop.asar", "chunma-san/ChunCord", directory); err == nil {
		t.Fatal("Corrupt download was accepted")
	}
	previous, err := os.ReadFile(destination)
	if err != nil || string(previous) != "previous" {
		t.Fatal("Existing installation changed")
	}
	files, _ := os.ReadDir(directory)
	if len(files) != 1 {
		t.Fatal("Failed download left temporary files")
	}
	corrupt = false
	temporary, err := downloadVerifiedArtifact(release, "desktop.asar", "chunma-san/ChunCord", directory)
	if err != nil {
		t.Fatal(err)
	}
	if err = replaceVerifiedArtifact(temporary, destination); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(destination)
	if err != nil || string(actual) != data {
		t.Fatal("Verified replacement failed")
	}
	if _, _, err = artifactURL(release, "desktop.asar", "other/repository"); err == nil {
		t.Fatal("Foreign repository was accepted")
	}
	if _, err = expectedChecksum(fmt.Sprintf("%s  desktop.asar\n%s  desktop.asar\n", checksum, checksum), "desktop.asar"); err == nil {
		t.Fatal("Duplicate checksum was accepted")
	}
}
