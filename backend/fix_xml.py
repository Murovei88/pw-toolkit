import re

input_file = 'data/pw/samples/my_build_utf8.xml'
output_file = 'data/pw/samples/my_build_fixed.xml'

with open(input_file, 'r', encoding='utf-8') as f:
    xml_content = f.read()

# Убираем запятые и пробелы перед =" (например: SoulInfo00,="8" -> SoulInfo00="8")
fixed_content = re.sub(r',\s*="', '="', xml_content)

with open(output_file, 'w', encoding='utf-8') as f:
    f.write(fixed_content)

print(f"✅ XML исправлен и сохранен в {output_file}")
print("Теперь запусти парсер:")
print(f"go run cmd/data-parser/main.go {output_file}")
