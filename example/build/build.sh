#!/bin/sh

platform="$1"
extension="$2"

project_dir="$(realpath "$COMMAND_PACK_DIR/../../")" || {
    echo "failed to resolve COMMAND_PACK_DIR=$COMMAND_PACK_DIR/../../"
    exit 1
}
binaries_dir="$project_dir/.bin"
file_prefix="tuiPack"

if [ -n "$extension" ]
then
file_name="$file_prefix.$extension"
else
file_name="$file_prefix"
fi

binary_path="$binaries_dir/$file_name"
archive_name="$file_prefix-$platform.tar.gz"
archive_path="$binaries_dir/$archive_name"

rm -f "$archive_path"
rm -f "$binary_path"

(cd "$project_dir" && go build -o "$binary_path") || {
    echo "failed to build tuiPack binary"
    exit 1
}

(cd "$binaries_dir" && tar -czvf "$archive_path" "$file_name") || {
    rm -f "$binary_path"
    echo "failed to archive"
    exit 1
}

rm -f "$binary_path"
echo "archive created $archive_path"
