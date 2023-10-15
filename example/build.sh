#!/bin/sh

FILENAME="$1"
FILEPATH="$TUI_PACK_CONFIG_DIR/bin/$FILENAME"
rm -f "$FILEPATH"
(cd "$TUI_PACK_CONFIG_DIR" && go build -o "$FILEPATH") || {
    echo "failed to build"
    exit 1
}
(cd "$TUI_PACK_CONFIG_DIR/bin/" && tar -czvf "$FILENAME.tar.gz" "$FILENAME") || {
    echo "failed to archive"
    exit 1
}
