with open("src/core/db/postgres.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("userCache      *cache.Cache[\nUsers]", "userCache      *cache.Cache[*Users]")
c = c.replace("	\"github.com/jackc/pgx/v5\"\n", "")
with open("src/core/db/postgres.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("src/core/db/chats.go", "r", encoding="utf-8") as f:
    c = f.read()
if "\"time\"" not in c:
    c = c.replace("\"log/slog\"", "\"log/slog\"\n\t\"time\"")
with open("src/core/db/chats.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("src/core/db/auth.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("	\"github.com/jackc/pgx/v5\"\n", "")
with open("src/core/db/auth.go", "w", encoding="utf-8") as f:
    f.write(c)

with open("src/core/db/users.go", "r", encoding="utf-8") as f:
    c = f.read()
c = c.replace("db.userCache.Set(key, Users{ID: id})", "db.userCache.Set(key, &Users{ID: id})")
with open("src/core/db/users.go", "w", encoding="utf-8") as f:
    f.write(c)

