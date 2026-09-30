with open("src/vc/userbot.go", "r", encoding="utf-8") as f:
    text = f.read()

text = text.replace(
    'case strings.Contains(errMsg, "CHANNEL_PRIVATE"):',
    'case strings.Contains(errMsg, "PROMOTE ME OR ADD USERACCOUNT"):\n\t\treturn fmt.Errorf("PROMOTE ME OR ADD USERACCOUNT")\n\n\tcase strings.Contains(errMsg, "CHANNEL_PRIVATE"):'
)

with open("src/vc/userbot.go", "w", encoding="utf-8") as f:
    f.write(text)
