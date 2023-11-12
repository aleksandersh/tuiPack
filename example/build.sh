#!/bin/sh

platform="$1"
extension="$2"

project_dir="$(realpath "$TUI_PACK_CONFIG_DIR/../")"
binaries_dir="$project_dir/.bin"

if [ -n "$extension" ]
then
file_name="tuiPack.$extension"
else
file_name="tuiPack"
fi

binary_path="$binaries_dir/$file_name"
archive_name="$file_name-$platform.tar.gz"
archive_path="$binaries_dir/$archive_name"

rm -f "$archive_path"
rm -f "$binary_path"

(cd "$project_dir" && go build -o "$binary_path") || {
    echo "failed to build tuiPack binary"
    exit 1
}

(cd "$project_dir" && tar -czvf "$archive_path" "$file_name") || {
    rm -f "$binary_path"
    echo "failed to archive"
    exit 1
}

rm -f "$binary_path"
