#!/usr/bin/env bash
# Package SkillUI Linux release: deb + AppImage
#
# Usage: scripts/package-linux.sh <goarch> <version>
#   goarch:  amd64 | arm64
#   version: e.g. 0.3.0-beta
#
# Requires: wails build output at build/bin/SkillUI, build/appicon.png
# Outputs:  SkillUI-<version>-linux-<goarch>.deb / .AppImage (repo root)

set -euo pipefail

GOARCH="${1:?goarch required (amd64|arm64)}"
VERSION="${2:?version required}"

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

case "$GOARCH" in
  amd64)
    DEB_ARCH="amd64"
    APPIMAGE_TOOL="appimagetool-x86_64.AppImage"
    ;;
  arm64)
    DEB_ARCH="arm64"
    APPIMAGE_TOOL="appimagetool-aarch64.AppImage"
    ;;
  *)
    echo "ERROR: unsupported goarch: $GOARCH" >&2
    exit 1
    ;;
esac

DESKTOP_CONTENT='[Desktop Entry]
Name=SkillUI
Comment=Manage AI coding assistant skills easily
Exec=skillui
Icon=skillui
Terminal=false
Type=Application
Categories=Utility;Development;
'

# ── deb (via fpm) ──────────────────────────────────────────────
if ! command -v fpm >/dev/null 2>&1; then
  echo "Installing fpm..."
  sudo apt-get install -y ruby ruby-dev rubygems >/dev/null 2>&1
  sudo gem install fpm >/dev/null
fi

PKG_DIR="pkg-deb"
rm -rf "$PKG_DIR"
mkdir -p "$PKG_DIR/usr/bin" "$PKG_DIR/usr/share/applications" "$PKG_DIR/usr/share/icons/hicolor/256x256/apps"
cp build/bin/SkillUI "$PKG_DIR/usr/bin/skillui"
printf '%s' "$DESKTOP_CONTENT" > "$PKG_DIR/usr/share/applications/skillui.desktop"
cp build/appicon.png "$PKG_DIR/usr/share/icons/hicolor/256x256/apps/skillui.png"

fpm -s dir -t deb \
  -n skillui \
  -v "$VERSION" \
  -a "$DEB_ARCH" \
  --description "SkillUI - Manage AI coding assistant skills easily" \
  --category "Utility" \
  --maintainer "MZ <modstart@163.com>" \
  --url "https://skillui.com" \
  --license "Apache-2.0" \
  -C "$PKG_DIR" \
  -p "SkillUI-${VERSION}-linux-${GOARCH}.deb" >/dev/null
rm -rf "$PKG_DIR"
echo "✅ deb: SkillUI-${VERSION}-linux-${GOARCH}.deb"

# ── AppImage (via appimagetool) ────────────────────────────────
if [ ! -f appimagetool ]; then
  echo "Downloading appimagetool..."
  wget -q "https://github.com/AppImage/appimagetool/releases/download/continuous/${APPIMAGE_TOOL}" -O appimagetool
  chmod +x appimagetool
fi

APPDIR="AppDir"
rm -rf "$APPDIR"
mkdir -p "$APPDIR/usr/bin"
cp build/bin/SkillUI "$APPDIR/usr/bin/skillui"
cp build/appicon.png "$APPDIR/skillui.png"
printf '%s' "$DESKTOP_CONTENT" > "$APPDIR/skillui.desktop"
printf '#!/bin/sh\nexec "$(dirname "$0")/usr/bin/skillui" "$@"\n' > "$APPDIR/AppRun"
chmod +x "$APPDIR/AppRun"

APPIMAGE_EXTRACT_AND_RUN=1 ./appimagetool "$APPDIR" "SkillUI-${VERSION}-linux-${GOARCH}.AppImage" >/dev/null 2>&1
rm -rf "$APPDIR"
echo "✅ AppImage: SkillUI-${VERSION}-linux-${GOARCH}.AppImage"

echo ""
echo "Packaged files:"
ls -lh SkillUI-${VERSION}-linux-${GOARCH}.deb SkillUI-${VERSION}-linux-${GOARCH}.AppImage
