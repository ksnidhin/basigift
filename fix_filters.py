import re
with open("src/handlers/filters.go", "r", encoding="utf-8") as f:
    text = f.read()

pattern = re.compile(r'func playMode\(c \*td\.Client, m \*td\.Message\) bool \{\n\s*if m\.IsPrivate\(\) \{\n\s*return false\n\s*\}\n\n\s*chatID := m\.ChatID\(\)\n\n\s*return false\n\s*\}\n')
text = text.replace('chatID := m.ChatID()\n\n\t\treturn false\n\t}\n', 'chatID := m.ChatID()\n')

with open("src/handlers/filters.go", "w", encoding="utf-8") as f:
    f.write(text)
