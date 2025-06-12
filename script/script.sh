#!/bin/bash

# Проверка наличия аргумента
if [ -z "$1" ]; then
  echo "Укажи путь к папке с .txt файлами"
  echo "Пример: ./merge_txt.sh ./input"
  exit 1
fi

INPUT_DIR="$1"
OUTPUT_FILE="output-concat.txt"

# Удалим старый выходной файл, если существует
rm -f "$OUTPUT_FILE"

# Объединение всех .txt файлов
for file in "$INPUT_DIR"/*.txt; do
  echo "Добавляю: $file"
  cat "$file" >> "$OUTPUT_FILE"
  echo "" >> "$OUTPUT_FILE" # Добавим перенос строки между файлами
done

echo "Готово. Все файлы объединены в: $OUTPUT_FILE"
