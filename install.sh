#!/bin/bash
set -e

echo "======================================"
echo " Daddy Noah Native Ubuntu Installer"
echo "======================================"

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (use sudo)"
  exit 1
fi

echo "[1/6] Installing dependencies..."
apt-get update
apt-get install -y postgresql postgresql-contrib ffmpeg wget curl git build-essential

echo "[2/6] Installing yt-dlp..."
curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
chmod a+rx /usr/local/bin/yt-dlp

echo "[3/6] Installing Go (if needed)..."
if ! command -v go &> /dev/null; then
    GO_VERSION="1.22.4"
    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    rm -rf /usr/local/go && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    rm go${GO_VERSION}.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
fi

echo "[4/6] Setting up PostgreSQL..."
sudo -u postgres psql -c "CREATE USER postgres WITH SUPERUSER PASSWORD 'postgres';" || true
sudo -u postgres psql -c "CREATE DATABASE daddynoah;" || true

echo "[5/6] Building Daddy Noah..."
export PATH=$PATH:/usr/local/go/bin
mkdir -p /opt/DaddyNoah
cd /opt/DaddyNoah
# Assuming the user cloned the repo here, or we copy it:
cp -r $PWD/* /opt/DaddyNoah/ || true
go mod tidy
go build -o daddynoah .

echo "[6/6] Setting up systemd service..."
cp daddynoah.service /etc/systemd/system/daddynoah.service
systemctl daemon-reload
systemctl enable daddynoah

echo "======================================"
echo "Installation complete!"
echo ""
echo "Please configure your bot before starting:"
echo "1. cp sample.env .env"
echo "2. nano .env"
echo "3. systemctl start daddynoah"
echo "======================================"
