/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"equilotl/buildinfo"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

var IsSelfOutdated = false
var SelfUpdateCheckDoneChan = make(chan bool, 1)

func init() {
	//goland:noinspection GoBoolExpressions
	if buildinfo.InstallerTag == buildinfo.VersionUnknown {
		Log.Debug("Disabling self updater as this is not a release build")
		return
	}

	go DeleteOldExecutable()

	go func() {
		Log.Debug("Checking for Installer Updates...")

		res, err := GetGithubRelease(InstallerReleaseUrl, InstallerReleaseUrlFallback)
		if err != nil {
			Log.Warn("Failed to check for self updates:", err)
			SelfUpdateCheckDoneChan <- false
		} else {
			IsSelfOutdated = res.TagName != buildinfo.InstallerTag
			Log.Debug("Is self outdated?", IsSelfOutdated)
			SelfUpdateCheckDoneChan <- true
		}
	}()
}

func GetInstallerDownloadLink() string {
	const BaseUrl = "https://github.com/chunma-san/ChunCord-Installer/releases/latest/download/"
	switch runtime.GOOS {
	case "windows":
		filename := Ternary(buildinfo.UiType == buildinfo.UiTypeCli, "ChunCord-Installer-Cli", "ChunCord-Installer")
		if runtime.GOARCH == "arm64" {
			filename += "-arm64"
		}
		return BaseUrl + filename + ".exe"
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return BaseUrl + "ChunCord-Installer-x64.dmg"
		case "arm64":
			return BaseUrl + "ChunCord-Installer-arm64.dmg"
		default:
			return ""
		}
	case "linux":
		if runtime.GOARCH == "arm64" {
			return BaseUrl + "ChunCord-Installer-Cli-linux-arm64"
		}
		return BaseUrl + "ChunCord-Installer-Cli-linux"
	default:
		return ""
	}
}

func CanUpdateSelf() bool {
	//goland:noinspection GoBoolExpressions
	return IsSelfOutdated && runtime.GOOS == "windows"
}

func UpdateSelf() error {
	if !CanUpdateSelf() {
		return errors.New("Cannot update self. Either no update available or macos")
	}

	url := GetInstallerDownloadLink()
	if url == "" {
		return errors.New("Failed to get installer download link")
	}

	Log.Debug("Updating self from", url)

	ownExePath, err := os.Executable()
	if err != nil {
		return err
	}

	ownExeDir := path.Dir(ownExePath)

	release, err := GetGithubRelease(InstallerReleaseUrl, InstallerReleaseUrlFallback)
	if err != nil {
		return err
	}
	filename := path.Base(url)
	tmpName, err := downloadVerifiedArtifact(release, filename, "chunma-san/ChunCord-Installer", ownExeDir)
	if err != nil {
		return err
	}
	defer os.Remove(tmpName)

	if err = os.Rename(ownExePath, ownExePath+".old"); err != nil {
		return fmt.Errorf("Failed to preserve previous installer: %w", err)
	}

	if err = os.Rename(tmpName, ownExePath); err != nil {
		if restore := os.Rename(ownExePath+".old", ownExePath); restore != nil {
			return fmt.Errorf("Installer update failed: %v; restore failed: %v", err, restore)
		}
		return fmt.Errorf("Failed to replace self with updated executable. Please manually redownload the installer: %w", err)
	}

	return nil
}

func DeleteOldExecutable() {
	ownExePath, err := os.Executable()
	if err != nil {
		return
	}

	for attempts := 0; attempts < 10; attempts += 1 {
		err = os.Remove(ownExePath + ".old")

		if err == nil || errors.Is(err, os.ErrNotExist) {
			break
		}

		Log.Warn("Failed to remove old executable. Retrying in 1 second.", err)
		time.Sleep(1 * time.Second)
	}
}

func RelaunchSelf() error {
	attr := new(os.ProcAttr)
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	var argv []string
	if len(os.Args) > 1 {
		argv = os.Args[1:]
	} else {
		argv = []string{}
	}

	Log.Debug("Restarting self with exe", os.Args[0], "and args", argv)

	proc, err := os.StartProcess(os.Args[0], argv, attr)
	if err != nil {
		return fmt.Errorf("Failed to start new process: %w", err)
	}

	if err = proc.Release(); err != nil {
		return fmt.Errorf("Failed to release new process: %w", err)
	}

	os.Exit(0)
	return nil
}
