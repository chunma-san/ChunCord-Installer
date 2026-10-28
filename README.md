<!-- CAT-COUNTER:START -->
<p align="center">
  <img src="https://count.himiku.com/@chunma-chuncord-installer?theme=capoo-2&amp;padding=6&amp;scale=0.8&amp;darkmode=1" alt="Cute cat project visits counter" />
</p>
<!-- CAT-COUNTER:END -->

# ChunCord Installer

Install, repair or remove [ChunCord](https://github.com/chunma-san/ChunCord) on Windows.

## Install ChunCord

1. Download **[ChunCord-Installer.exe](https://github.com/chunma-san/ChunCord-Installer/releases/latest/download/ChunCord-Installer.exe)**.
2. Close Discord and open the installer.
3. Select your Discord installation and click **Install**.
4. Start Discord and open **Settings → ChunCord → Plugins**.

No Git, Node.js or source folder is needed. The installer downloads ChunCord, checks the release checksums and places its files in your application data folder. You can delete the downloaded installer when setup is finished.

The current release supports **Windows x64**. It is unsigned, so Windows may show an unknown publisher warning. Download it from this repository's releases.

## Updates and repairs

ChunCord updates from **Settings → ChunCord → Updater** inside Discord. Updates are downloaded and verified in the background, then applied before ChunCord loads on the next start. Enable update notifications if you want a reminder to restart.

If a Discord update removes the modification, close Discord and run the latest installer again:

- **Repair** restores ChunCord while keeping your settings.
- **Uninstall** removes ChunCord from the selected Discord installation.

Older installations that show an `EPERM` update error need one repair with the latest installer. Existing Equicord settings paths remain compatible.

Client downloads and installer updates use the ChunCord repositories. The installer does not fall back to an official Equicord build.

Installed-version checks use the release's source commit metadata, with support for legacy revision suffixes in release titles. A clean display name such as `ChunCord` does not replace its source revision.

## Other downloads

The [latest release](https://github.com/chunma-san/ChunCord-Installer/releases/latest) also includes a command-line installer and SHA-256 checksums. Most users only need **ChunCord-Installer.exe**.

## Building from source

Use the Go version specified in `go.mod`. Building the Windows GUI also requires MinGW-w64 and SDL2. The release workflow builds and tests the GUI and command-line versions before publishing them.

## Maintainer and credits

Maintained by [chunma-san](https://github.com/chunma-san).

Based on [Equilotl](https://github.com/Equicord/Equilotl) and the [Vencord Installer](https://github.com/Vencord/Installer). Thanks to their original contributors. Source headers and the [GPL license](LICENSE) are preserved.

ChunCord Installer is an independent fork and is not affiliated with Discord, Equicord or Vencord.
