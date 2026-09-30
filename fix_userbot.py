import re
with open("src/vc/userbot.go", "r", encoding="utf-8") as f:
    text = f.read()

text = re.sub(
    r'if err != nil \{\s*return "", fmt\.Errorf\("create invite link for chat %d: %w", chatID, err\)\s*\}',
    'if err != nil {\n\t\treturn "", fmt.Errorf("PROMOTE ME OR ADD USERACCOUNT")\n\t}',
    text
)

with open("src/vc/userbot.go", "w", encoding="utf-8") as f:
    f.write(text)
