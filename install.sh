#!/bin/bash
set -e

REPO_DIR="$PWD"

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
if ! command -v go &> /dev/null || ! /usr/local/go/bin/go version &> /dev/null; then
    GO_VERSION="1.22.4"
    ARCH=$(uname -m)
    if [ "$ARCH" = "x86_64" ]; then
        GO_ARCH="amd64"
    elif [ "$ARCH" = "aarch64" ]; then
        GO_ARCH="arm64"
    else
        echo "Unsupported architecture: $ARCH"
        exit 1
    fi
    wget "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    rm -rf /usr/local/go && tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    rm "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    export PATH=$PATH:/usr/local/go/bin
fi

echo "[4/6] Setting up PostgreSQL..."
sudo -u postgres psql -c "CREATE USER postgres WITH SUPERUSER PASSWORD 'postgres';" || true
sudo -u postgres psql -c "CREATE DATABASE daddynoah;" || true

echo "[5/6] Building Daddy Noah..."
export PATH=$PATH:/usr/local/go/bin
mkdir -p /opt/DaddyNoah

# Copy everything to the deployment directory
cp -a "$REPO_DIR"/. /opt/DaddyNoah/
cd /opt/DaddyNoah

go mod tidy
go run setup_ntgcalls.go
go build -o daddynoah .

echo "[6/6] Setting up systemd service..."
cp daddynoah.service /etc/systemd/system/daddynoah.service
systemctl daemon-reload
systemctl enable daddynoah

echo "======================================"
echo "Installation complete!"
echo ""
echo "Please configure your bot before starting:"
echo "1. cd /opt/DaddyNoah"
echo "2. cp sample.env .env"
echo "3. nano .env"
echo "4. sudo systemctl start daddynoah"
echo "======================================"
