#!/bin/bash
# ============================================================
#  AI Slayd Bot — Server Setup Script
#  Server ichidan ishga tushirish: bash deploy.sh
# ============================================================

set -e

# ── Ranglar ──────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'
info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC}   $1"; }
error()   { echo -e "${RED}[ERR]${NC}  $1"; exit 1; }

APP_DIR="/opt/ai-slayd-bot"
SERVICE_NAME="ai-slayd-bot"
GO_VERSION="1.23.4"

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  🚀  AI Slayd Bot — Server Setup${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# ── Go o'rnatish / yangilash ──────────────────────────────────
INSTALLED_GO=""
if [ -f /usr/local/go/bin/go ]; then
    INSTALLED_GO=$(/usr/local/go/bin/go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
fi

if [ "$INSTALLED_GO" = "$GO_VERSION" ]; then
    success "Go $GO_VERSION allaqachon o'rnatilgan"
else
    info "Go $GO_VERSION o'rnatilmoqda (mavjud: ${INSTALLED_GO:-yo'q})..."
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    success "Go $GO_VERSION o'rnatildi"
fi

export PATH=$PATH:/usr/local/go/bin
grep -qxF 'export PATH=$PATH:/usr/local/go/bin' ~/.bashrc || echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
info "Go versiyasi: $(/usr/local/go/bin/go version)"

# ── .env fayl yaratish ───────────────────────────────────────
info ".env fayli yaratilmoqda..."
cat > "$APP_DIR/.env" << 'ENVEOF'
TELEGRAM_BOT_TOKEN=8612614565:AAFdhvhqyKAFtVCc105M2ebASUXc8Nakpqc
DEEPSEEK_API_KEY=sk-2c84e1dd866c4d0a9c102cdca1e719f6
TMP_DIR=./tmp
ENVEOF
success ".env fayli yaratildi"

# ── tmp papkasi ───────────────────────────────────────────────
mkdir -p "$APP_DIR/tmp"

# ── Build qilish ─────────────────────────────────────────────
info "Bot build qilinmoqda..."
cd "$APP_DIR"
/usr/local/go/bin/go mod download
/usr/local/go/bin/go build -o ai-slayd-bot .
success "Build muvaffaqiyatli!"

# ── Systemd service ───────────────────────────────────────────
info "Systemd service yaratilmoqda..."
cat > /etc/systemd/system/${SERVICE_NAME}.service << SVCEOF
[Unit]
Description=AI Slayd Bot - Telegram Presentation Generator
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${APP_DIR}
ExecStart=${APP_DIR}/ai-slayd-bot
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
EnvironmentFile=${APP_DIR}/.env

[Install]
WantedBy=multi-user.target
SVCEOF

systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
systemctl restart "$SERVICE_NAME"
sleep 3

# ── Natija ───────────────────────────────────────────────────
echo ""
if systemctl is-active --quiet "$SERVICE_NAME"; then
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}  ✅  Bot muvaffaqiyatli ishga tushdi!${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "  📊 Status:  ${YELLOW}systemctl status $SERVICE_NAME${NC}"
    echo -e "  📋 Loglar:  ${YELLOW}journalctl -u $SERVICE_NAME -f${NC}"
    echo -e "  🔄 Restart: ${YELLOW}systemctl restart $SERVICE_NAME${NC}"
    echo -e "  🛑 Stop:    ${YELLOW}systemctl stop $SERVICE_NAME${NC}"
else
    echo -e "${RED}❌ Bot ishga tushmadi. Loglar:${NC}"
    journalctl -u "$SERVICE_NAME" -n 30 --no-pager
    exit 1
fi
echo ""