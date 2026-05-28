#!/usr/bin/env bash
set -euo pipefail

APP_NAME="Headlamp Agent"
BUNDLE_ID="cloud.unboxd.headlamp-agent"
DIST_DIR="dist"
APP_DIR="$DIST_DIR/$APP_NAME.app"
CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
RESOURCES_DIR="$CONTENTS_DIR/Resources"
DMG_PATH="$DIST_DIR/Headlamp-Agent.dmg"

mkdir -p "$MACOS_DIR" "$RESOURCES_DIR"

go build -o "$MACOS_DIR/headlamp" ./cmd/headlamp
go build -o "$MACOS_DIR/headlamp-node-agent" ./cmd/headlamp-node-agent

cat > "$MACOS_DIR/HeadlampAgent" <<'SH'
#!/usr/bin/env bash
DIR="$(cd "$(dirname "$0")" && pwd)"
open -a Terminal "$DIR/headlamp"
SH
chmod +x "$MACOS_DIR/HeadlampAgent"

cat > "$CONTENTS_DIR/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleExecutable</key>
  <string>HeadlampAgent</string>
  <key>CFBundleIdentifier</key>
  <string>$BUNDLE_ID</string>
  <key>CFBundleName</key>
  <string>$APP_NAME</string>
  <key>CFBundleDisplayName</key>
  <string>$APP_NAME</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>CFBundleShortVersionString</key>
  <string>0.1.0</string>
  <key>CFBundleVersion</key>
  <string>0.1.0</string>
</dict>
</plist>
PLIST

rm -f "$DMG_PATH"
hdiutil create -volname "$APP_NAME" -srcfolder "$APP_DIR" -ov -format UDZO "$DMG_PATH"

echo "Built $DMG_PATH"
