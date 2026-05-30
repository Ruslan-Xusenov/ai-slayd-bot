set -e

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'
info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC}   $1"; }
error()   { echo -e "${RED}[ERR]${NC}  $1"; exit 1; }

APP_DIR="/opt/ai-slayd-bot"
SERVICE_NAME="ai-slayd-bot"
GO_VERSION="1.23.4"

echo ""
echo -e "${BLUE}===================================================${NC}"
echo -e "${BLUE}  AI Slayd Bot - Server Setup${NC}"
echo -e "${BLUE}===================================================${NC}"
echo ""

INSTALLED_GO=""
if [ -f /usr/local/go/bin/go ]; then
    INSTALLED_GO=$(/usr/local/go/bin/go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
fi

if [ "$INSTALLED_GO" = "$GO_VERSION" ]; then
    success "Go $GO_VERSION allaqachon mavjud"
else
    info "Go $GO_VERSION yuklanmoqda (mavjud: ${INSTALLED_GO:-yoq})..."
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    success "Go $GO_VERSION o'rnatildi"
fi

export PATH=$PATH:/usr/local/go/bin
grep -qF '/usr/local/go/bin' ~/.bashrc || echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
info "Go: $(/usr/local/go/bin/go version)"

info ".env fayli yaratilmoqda..."
printf 'TELEGRAM_BOT_TOKEN=8612614565:AAFdhvhqyKAFtVCc105M2ebASUXc8Nakpqc\n' > "$APP_DIR/.env"
printf 'DEEPSEEK_API_KEY=sk-2c84e1dd866c4d0a9c102cdca1e719f6\n' >> "$APP_DIR/.env"
printf 'TMP_DIR=./tmp\n' >> "$APP_DIR/.env"
success ".env yaratildi"

mkdir -p "$APP_DIR/tmp"

info "Bot build qilinmoqda..."
cd "$APP_DIR"
/usr/local/go/bin/go mod download
/usr/local/go/bin/go build -o ai-slayd-bot .
success "Build OK!"

info "Systemd service yaratilmoqda..."
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

printf '[Unit]\n' > "$SERVICE_FILE"
printf 'Description=AI Slayd Bot - Telegram Presentation Generator\n' >> "$SERVICE_FILE"
printf 'After=network.target\n' >> "$SERVICE_FILE"
printf '\n' >> "$SERVICE_FILE"
printf '[Service]\n' >> "$SERVICE_FILE"
printf 'Type=simple\n' >> "$SERVICE_FILE"
printf 'User=root\n' >> "$SERVICE_FILE"
printf "WorkingDirectory=${APP_DIR}\n" >> "$SERVICE_FILE"
printf "ExecStart=${APP_DIR}/ai-slayd-bot\n" >> "$SERVICE_FILE"
printf 'Restart=always\n' >> "$SERVICE_FILE"
printf 'RestartSec=5\n' >> "$SERVICE_FILE"
printf 'StandardOutput=journal\n' >> "$SERVICE_FILE"
printf 'StandardError=journal\n' >> "$SERVICE_FILE"
printf "EnvironmentFile=${APP_DIR}/.env\n" >> "$SERVICE_FILE"
printf '\n' >> "$SERVICE_FILE"
printf '[Install]\n' >> "$SERVICE_FILE"
printf 'WantedBy=multi-user.target\n' >> "$SERVICE_FILE"

systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
systemctl restart "$SERVICE_NAME"
sleep 3

echo ""
if systemctl is-active --quiet "$SERVICE_NAME"; then
    echo -e "${GREEN}===================================================${NC}"
    echo -e "${GREEN}  Bot muvaffaqiyatli ishga tushdi!${NC}"
    echo -e "${GREEN}===================================================${NC}"
    echo ""
    echo -e "  Loglar:   ${YELLOW}journalctl -u $SERVICE_NAME -f${NC}"
    echo -e "  Status:   ${YELLOW}systemctl status $SERVICE_NAME${NC}"
    echo -e "  Restart:  ${YELLOW}systemctl restart $SERVICE_NAME${NC}"
    echo -e "  Stop:     ${YELLOW}systemctl stop $SERVICE_NAME${NC}"
else
    echo -e "${RED}Bot ishga tushmadi. Loglar:${NC}"
    journalctl -u "$SERVICE_NAME" -n 30 --no-pager
    exit 1
fi
echo ""