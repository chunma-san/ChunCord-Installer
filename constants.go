/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"equilotl/buildinfo"
	"image/color"
)

const ReleaseUrl = "https://api.github.com/repos/chunma-san/ChunCord/releases/latest"
const ReleaseUrlFallback = ReleaseUrl
const InstallerReleaseUrl = "https://api.github.com/repos/chunma-san/ChunCord-Installer/releases/latest"
const InstallerReleaseUrlFallback = InstallerReleaseUrl

var UserAgent = "ChunCord-Installer/" + buildinfo.InstallerGitHash + " (https://github.com/chunma-san/ChunCord-Installer)"

const SupportUrl = "https://github.com/chunma-san/ChunCord"

var (
	DiscordGreen        = color.RGBA{0, 133, 69, 0xff}
	DiscordGreenHovered = color.RGBA{0, 108, 55, 0xff}
	DiscordRed          = color.RGBA{210, 45, 57, 0xff}
	DiscordRedHovered   = color.RGBA{169, 35, 46, 0xff}
	DiscordBlue         = color.RGBA{88, 101, 242, 0xff}
	DiscordBlueHovered  = color.RGBA{68, 82, 187, 0xff}
	DiscordYellow       = color.RGBA{0xfe, 0xe7, 0x5c, 0xff}
)

var LinuxDiscordNames = []string{
	"Discord",
	"DiscordPTB",
	"DiscordCanary",
	"DiscordDevelopment",
	"discord",
	"discordptb",
	"discordcanary",
	"discorddevelopment",
	"discord-ptb",
	"discord-canary",
	"discord-development",
	// Flatpak
	"com.discordapp.Discord",
	"com.discordapp.DiscordPTB",
	"com.discordapp.DiscordCanary",
	"com.discordapp.DiscordDevelopment",
}
