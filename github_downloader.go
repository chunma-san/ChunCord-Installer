/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	path "path/filepath"
	"regexp"
	"strings"
)

type GithubRelease struct {
	Name            string `json:"name"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Assets          []struct {
		Name        string `json:"name"`
		Size        int64  `json:"size"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var ReleaseData GithubRelease
var GithubError error
var GithubDoneChan chan bool

var InstalledHash = "None"
var LatestHash = "Unknown"
var IsDevInstall bool

// latestRevision returns the source revision of a release. Releases are titled
// "ChunCord" without any appended hash, so the revision comes from
// target_commitish (set by the release workflow); the title suffix is kept as a
// fallback for releases that still carry one.
func latestRevision(release *GithubRelease) string {
	if match := regexp.MustCompile(`^[a-f0-9]{7,40}$`).FindString(strings.ToLower(release.TargetCommitish)); match != "" {
		return match
	}
	i := strings.LastIndex(release.Name, " ") + 1
	return release.Name[i:]
}

func GetGithubRelease(url, fallbackUrl string) (*GithubRelease, error) {
	Log.Debug("Fetching", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error("Failed to create Request", err)
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Log.Error("Failed to send Request", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		isRateLimitedOrBlocked := res.StatusCode == 401 || res.StatusCode == 403 || res.StatusCode == 429
		triedFallback := url == fallbackUrl

		if isRateLimitedOrBlocked && !triedFallback {
			Log.Error(fmt.Sprintf("Failed to fetch %s (status code %d). Trying fallback url %s", url, res.StatusCode, fallbackUrl))
			return GetGithubRelease(fallbackUrl, fallbackUrl)
		}

		err = errors.New(res.Status)
		Log.Error(url, "returned Non-OK status", GithubError)
		return nil, err
	}

	var data GithubRelease

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		Log.Error("Failed to decode GitHub JSON Response", err)
		return nil, err
	}

	return &data, nil
}

func InitGithubDownloader() {
	GithubDoneChan = make(chan bool, 1)

	IsDevInstall = os.Getenv("EQUICORD_DEV_INSTALL") == "1"
	Log.Debug("Is Dev Install: ", IsDevInstall)
	if IsDevInstall {
		GithubDoneChan <- true
		return
	}

	go func() {
		// Make sure UI updates once the request either finished or failed
		defer func() {
			GithubDoneChan <- GithubError == nil
		}()

		data, err := GetGithubRelease(ReleaseUrl, ReleaseUrlFallback)
		if err != nil {
			GithubError = err
			return
		}

		ReleaseData = *data

		LatestHash = latestRevision(data)
		Log.Debug("Finished fetching GitHub Data")
		Log.Debug("Latest hash is", LatestHash, "Local Install is", Ternary(LatestHash == InstalledHash, "up to date!", "outdated!"))
	}()

	// either .asar file or directory with main.js file (in DEV)
	EquicordFile := EquicordDirectory

	stat, err := os.Stat(EquicordFile)
	if err != nil {
		return
	}

	// dev
	if stat.IsDir() {
		EquicordFile = path.Join(EquicordFile, "main.js")
	}

	// Check hash of installed version if exists
	b, err := os.ReadFile(EquicordFile)
	if err != nil {
		return
	}

	Log.Debug("Found existing ChunCord Install. Checking for hash...")

	re := regexp.MustCompile(`// Equicord (\w+)`)
	match := re.FindSubmatch(b)
	if match != nil {
		InstalledHash = string(match[1])
		Log.Debug("Existing hash is", InstalledHash)

	} else {
		Log.Debug("Didn't find hash")

	}
}

func installLatestBuilds() (retErr error) {
	Log.Debug("Installing latest builds...")

	if IsDevInstall {
		Log.Debug("Skipping due to dev install")
		return
	}

	temporary, err := downloadVerifiedArtifact(&ReleaseData, "desktop.asar", "chunma-san/ChunCord", path.Dir(EquicordDirectory))
	if err != nil {
		return err
	}
	defer os.Remove(temporary)
	if err = replaceVerifiedArtifact(temporary, EquicordDirectory); err != nil {
		return err
	}
	_ = FixOwnership(EquicordDirectory)
	InstalledHash = LatestHash
	return nil
}
