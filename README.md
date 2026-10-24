# ChunCord Installer

A Windows installer for [ChunCord](https://github.com/chunma-san/ChunCord), based on [Equilotl](https://github.com/Equicord/Equilotl).

## Install

1. Download **[ChunCord-Installer.exe](https://github.com/chunma-san/ChunCord-Installer/releases/latest/download/ChunCord-Installer.exe)**.
2. Close Discord, open the installer, choose your Discord installation and click **Install**.
3. Start Discord and open **Settings → ChunCord → Plugins**.

No Git, Node.js or source folder is needed. The installer stores ChunCord in the user's application data folder. You can delete the downloaded installer after setup. Existing Equicord settings paths remain compatible.

ChunCord checks its own GitHub releases for updates inside Discord. Keep update notifications enabled to know when to restart. Run this installer again for **Repair** after a Discord update, or use **Uninstall** to restore Discord.

The published installer targets Windows x64. It is currently unsigned; Windows can display an unknown publisher prompt. Download it only from this repository's releases. A CLI build and SHA-256 checksums are included for users who need them.

## Updates

Client downloads, fallback handling and installer self-updates point to ChunCord repositories. They do not fall back to installing official Equicord builds. Downloads are checked against SHA-256 release checksums before replacing an installed file.

## Build

Use the Go version in `go.mod`. The Windows GUI also requires MinGW-w64 and SDL2. The release workflow builds the GUI and CLI, runs tests and publishes only successful builds. Published artifacts contain no local settings or account data.

## Contributors

- [chunma-san](https://github.com/chunma-san): ChunCord installer maintenance, integration and releases.
- Original [Equilotl contributors](https://github.com/Equicord/Equilotl/graphs/contributors) and [Vencord Installer contributors](https://github.com/Vencord/Installer/graphs/contributors).

Original source headers and the [GPL license](LICENSE) are preserved. This is an independent fork and is not affiliated with Discord, Equicord or Vencord.
