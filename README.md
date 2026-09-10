# Daddy Noah

A high-performance, native Telegram Music Bot written in Go, powered by PostgreSQL, FFmpeg, yt-dlp, and systemd.

## Features
- Plays music and video in Telegram voice chats.
- Supports YouTube, Spotify, and more.
- High reliability with automatic recovery from crashes and disconnects.
- Fully native deployment on Ubuntu (No Docker required).
- PostgreSQL data persistence.

## Requirements
- Ubuntu 20.04+ (or similar Linux distro)
- PostgreSQL
- FFmpeg
- yt-dlp
- Go 1.21+

## Quick Installation (Ubuntu VPS)

```bash
git clone https://github.com/Simmie/DaddyNoah
cd DaddyNoah
sudo bash install.sh
```

## Manual Native Ubuntu Installation

1. **Install Dependencies:**
```bash
sudo apt update
sudo apt install -y postgresql postgresql-contrib ffmpeg curl
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp
```

2. **PostgreSQL setup:**
```bash
sudo -u postgres psql -c "CREATE USER postgres WITH SUPERUSER PASSWORD 'postgres';"
sudo -u postgres psql -c "CREATE DATABASE daddynoah;"
```

3. **Environment configuration:**
```bash
cp sample.env .env
nano .env
```
Ensure you configure:
- `TOKEN` (Telegram Bot Token)
- `API_ID` & `API_HASH` (Telegram API credentials)
- `STRING1` (Pyrogram Session string for assistant)
- `DATABASE_URL` (e.g. `postgres://postgres:postgres@localhost:5432/daddynoah?sslmode=disable`)

4. **Build the binary:**
```bash
go build -o daddynoah .
```

## systemd Setup
To run the bot as a background service with automatic restarts:

1. Edit the provided `daddynoah.service` if your working directory differs from `/opt/DaddyNoah`.
2. Install the service:
```bash
sudo cp daddynoah.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable daddynoah
```

## Starting/stopping/restarting
```bash
sudo systemctl start daddynoah
sudo systemctl stop daddynoah
sudo systemctl restart daddynoah
```

## Logs
To view the live logs:
```bash
sudo journalctl -u daddynoah -f
```

## Updating
```bash
cd /opt/DaddyNoah
git pull
go build -o daddynoah .
sudo systemctl restart daddynoah
```

## Troubleshooting
- If you encounter PostgreSQL connection errors, ensure `DATABASE_URL` is correct.
- If songs fail to play, verify `ffmpeg` and `yt-dlp` are installed and working.

## License / required attribution
*  Copyright (c) 2025-2026 Ashok Shau
*  Licensed under GNU GPL v3
