#!/bin/sh

FILENAME="$1"
WORKING_DIR="$TUI_PACK_CONFIG_DIR/../"
FILEPATH="$WORKING_DIR/bin/$FILENAME"
rm -f "$FILEPATH"
(cd "$WORKING_DIR" && go build -o "$FILEPATH") || {
    echo "failed to build"
    exit 1
}
(cd "$WORKING_DIR/bin/" && tar -czvf "$FILENAME.tar.gz" "$FILENAME") || {
    echo "failed to archive"
    exit 1
}
