with open("src/vc/play.go", "r", encoding="utf-8") as f:
    text = f.read()

text = text.replace(
    'strings.Contains(msg, "GROUPCALL_FORBIDDEN"):',
    'strings.Contains(msg, "GROUPCALL_FORBIDDEN"),\n\t\tstrings.Contains(msg, "PROMOTE ME OR ADD USERACCOUNT"):'
)

text = text.replace(
    'if strings.Contains(msg, "GROUPCALL_INVALID") {',
    'if strings.Contains(msg, "PROMOTE ME OR ADD USERACCOUNT") {\n\t\treturn errors.New("PROMOTE ME OR ADD USERACCOUNT")\n\t}\n\n\tif strings.Contains(msg, "GROUPCALL_INVALID") {'
)

with open("src/vc/play.go", "w", encoding="utf-8") as f:
    f.write(text)
