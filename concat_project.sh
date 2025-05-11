#!/bin/bash
if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <directory> [output_file]"
  exit 1
fi

source_dir="$1"
output_file="${2:-output.txt}"

if [ ! -d "$source_dir" ]; then
  echo "Error: Directory '$source_dir' does not exist"
  exit 1
fi

# Create output_file if not exist
touch $output_file

# Get absolute paths
source_dir=$(realpath "$source_dir")
output_file=$(realpath "$output_file")

# Clear output file
> "$output_file"

# Add git log if CWD is a git repo
if [ -d ".git" ]; then
  echo "# Git Log (git adog3)" > "$output_file"
  git log --all --decorate --oneline --graph >> "$output_file"
  echo -e "\n# ----------------------------------------\n" >> "$output_file"
fi

# Find all files and concatenate them
find "$source_dir" -type f -not -path "*/\.*" | sort | while read -r file; do
  # Skip the output file itself
  if [ "$file" = "$output_file" ]; then
    continue
  fi
  
  # Skip binary files, object files, etc.
  mimetype=$(file --mime-type -b "$file")
  if [[ "$mimetype" != text/* ]]; then
    continue
  fi
  
  # Add file path as header (relative to source_dir)
  rel_path="${file#$source_dir/}"
  echo "# $rel_path" >> "$output_file"
  
  # Concatenate file content
  cat "$file" >> "$output_file"
  
  # Add newlines for separation
  echo -e "\n\n" >> "$output_file"
done

echo "Files concatenated to $output_file"
