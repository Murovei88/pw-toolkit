import re

file_path = 'internal/calc/calculator.go'
with open(file_path, 'r', encoding='utf-8') as f:
    content = f.read()

# 1. Добавляем импорт regexp, если его нет
if '"regexp"' not in content:
    content = content.replace('"encoding/xml"', '"encoding/xml"\n\t"regexp"')

# 2. Ищем вызов xml.Unmarshal и захватываем имя переменной с байтами
# Шаблон поймает: err := xml.Unmarshal(xmlData, &infos)
# Или: if err := xml.Unmarshal(data, &info); err != nil {
pattern = r'([ \t]*)(?:[a-zA-Z0-9_]+\s*:?=\s*)?xml\.Unmarshal\(\s*([a-zA-Z0-9_]+)\s*,'

def replacer(match):
    indent = match.group(1)
    var_name = match.group(2) # Имя переменной (например, xmlData или data)
    
    fix_code = (
        f'{indent}// Fix invalid XML attributes (e.g., SoulInfo00,="8" -> SoulInfo00="8")\n'
        f'{indent}re := regexp.MustCompile(`,\s*="`)\n'
        f'{indent}{var_name} = re.ReplaceAll({var_name}, []byte(`="`))\n'
        f'{indent}'
    )
    return fix_code + match.group(0)

if re.search(pattern, content):
    # count=1, чтобы применить только к первому найденному вызову
    content = re.sub(pattern, replacer, content, count=1)
    print("✅ Фикс успешно внедрен в Go код перед xml.Unmarshal!")
else:
    print("⚠️ Вызов xml.Unmarshal не найден. Проверьте путь к файлу или его содержимое.")

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)
