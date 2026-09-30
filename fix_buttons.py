with open("src/core/buttons.go", "r", encoding="utf-8") as f:
    c = f.read()

c = c.replace('var SourceCodeBtn = url("Source Code", "https://github.com/Simmie/DaddyNoah", gotdbot.ButtonStylePrimary{})', '')
c = c.replace('var channelBtn = url("Updates", config.SupportChannel, gotdbot.ButtonStyleDefault{})', '')
c = c.replace('var groupBtn = url("Group", config.SupportGroup, gotdbot.ButtonStyleDefault{})', 'var creditsBtn = url("Developer", "https://t.me/morevicodin", gotdbot.ButtonStyleDefault{})')

c = c.replace('{channelBtn, groupBtn}', '{creditsBtn}')
c = c.replace('{CloseBtn, SourceCodeBtn}', '{CloseBtn}')
c = c.replace('{SourceCodeBtn},', '')
c = c.replace('{SourceCodeBtn}', '')

with open("src/core/buttons.go", "w", encoding="utf-8") as f:
    f.write(c)
