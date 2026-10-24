/*
 * SPDX-License-Identifier: GPL-3.0
 * Based on Vencord Installer and Equilotl.
 * ChunCord integration maintained by chunma-san.
 */
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var downloadClient = &http.Client{Timeout: 3 * time.Minute}

func artifactURL(release *GithubRelease, filename, repository string) (string, int64, error) {
	prefix := "https://github.com/" + repository + "/releases/download/"
	for _, asset := range release.Assets {
		if asset.Name != filename {
			continue
		}
		if !strings.HasPrefix(asset.DownloadURL, prefix) || asset.Size <= 0 || asset.Size > 128*1024*1024 {
			return "", 0, errors.New("Invalid ChunCord download information")
		}
		return asset.DownloadURL, asset.Size, nil
	}
	return "", 0, fmt.Errorf("Release is missing %s", filename)
}

func expectedChecksum(text, filename string) (string, error) {
	if len(text) > 16384 {
		return "", errors.New("Checksum list is too large")
	}
	found := ""
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 || strings.TrimPrefix(fields[1], "*") != filename {
			continue
		}
		decoded, err := hex.DecodeString(fields[0])
		if err != nil || len(decoded) != sha256.Size || found != "" {
			return "", errors.New("Invalid or duplicate checksum")
		}
		found = strings.ToLower(fields[0])
	}
	if scanner.Err() != nil {
		return "", scanner.Err()
	}
	if found == "" {
		return "", errors.New("Download checksum was not found")
	}
	return found, nil
}

func downloadVerifiedArtifact(release *GithubRelease, filename, repository, directory string) (string, error) {
	assetURL, size, err := artifactURL(release, filename, repository)
	if err != nil {
		return "", err
	}
	checksumURL, _, err := artifactURL(release, "SHA256SUMS.txt", repository)
	if err != nil {
		return "", err
	}
	response, err := downloadClient.Get(checksumURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Checksum download failed: %s", response.Status)
	}
	list, err := io.ReadAll(io.LimitReader(response.Body, 16385))
	if err != nil {
		return "", err
	}
	expected, err := expectedChecksum(string(list), filename)
	if err != nil {
		return "", err
	}
	response, err = downloadClient.Get(assetURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Download failed: %s", response.Status)
	}
	file, err := os.CreateTemp(directory, ".chuncord-download-*")
	if err != nil {
		return "", err
	}
	name := file.Name()
	success := false
	defer func() {
		file.Close()
		if !success {
			os.Remove(name)
		}
	}()
	hash := sha256.New()
	received, err := io.Copy(io.MultiWriter(file, hash), io.LimitReader(response.Body, size+1))
	if err != nil {
		return "", err
	}
	if received != size || hex.EncodeToString(hash.Sum(nil)) != expected {
		return "", errors.New("ChunCord download verification failed; installed files were preserved")
	}
	if err = file.Sync(); err != nil {
		return "", err
	}
	if err = file.Chmod(0755); err != nil {
		return "", err
	}
	if err = file.Close(); err != nil {
		return "", err
	}
	success = true
	return name, nil
}

func replaceVerifiedArtifact(temporary, destination string) error {
	if err := os.Rename(temporary, destination); err == nil {
		return nil
	}
	backup := destination + ".chuncord-old"
	if _, err := os.Stat(backup); err == nil {
		return errors.New("A previous backup exists; repair the installation before retrying")
	}
	if err := os.Rename(destination, backup); err != nil {
		return err
	}
	if err := os.Rename(temporary, destination); err != nil {
		if restore := os.Rename(backup, destination); restore != nil {
			return fmt.Errorf("Update failed: %v; restore failed: %v", err, restore)
		}
		return err
	}
	return os.Remove(backup)
}
