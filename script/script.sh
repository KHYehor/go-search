#!/bin/bash

# Check an argumetn
if [ -z "$1" ]; then
  echo "Specify the path to a txt"
  echo "Example: ./merge_txt.sh ./input"
  exit 1
fi

INPUT_DIR="$1"
OUTPUT_FILE="output.txt"

# Delete the previous file
rm -f "$OUTPUT_FILE"

# Concat concated files
for i in {1..5}; do
  # Concat .txt files
  for file in "$INPUT_DIR"/*.txt; do
    echo "Concating: $file"
    cat "$file" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE" # Add separate line between
  done
done

echo "Finished, check the output: $OUTPUT_FILE"
