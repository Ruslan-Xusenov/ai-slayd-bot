#!/bin/bash
# ============================================================
#  AI Slayd Bot — Server Deploy Script
#  Ishlatish: bash deploy.sh <server_ip> <ssh_user> <ssh_password>
#  Misol:     bash deploy.sh 46.224.133.140 root mypassword
# ============================================================

set -e

# ── Argumentlar ──────────────────────────────────────────────
SERVER_IP="${1:-46.224.133.140}"
SSH_USER="${2:-root}"
SSH_PASS="${3}"
GITHUB_REPO="https://github.com/Ruslan-Xusenov/ai-slayd-bot.git"
APP_DIR="/opt/ai-slayd-bot"
SERVICE_NAME="ai-slayd-bot"

# ── Ranglar ──────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'
info()    { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC}   $1"; }
warn()    { echo -e "${YELLOW}[WARN]${NC} $1"; }
error()   { echo -e "${RED}[ERR]${NC}  $1"; exit 1; }

# ── Tekshirish ───────────────────────────────────────────────
if [ -z "$SSH_PASS" ]; then
    error "SSH parol kiritilmadi!\nIshlatish: bash deploy.sh 46.224.133.140 root YOUR_PASSWORD"
fi

command -v sshpass &>/dev/null || {
    warn "sshpass o'rnatilmagan. O'rnatilmoqda..."
    sudo apt-get install -y sshpass &>/dev/null
}

SSH_CMD="sshpass -p '$SSH_PASS' ssh -o StrictHostKeyChecking=no $SSH_USER@$SERVER_IP"
SCP_CMD="sshpass -p '$SSH_PASS' scp -o StrictHostKeyChecking=no"

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  🚀  AI Slayd Bot — Server Deploy${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# ── .env faylini tekshirish ───────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"

if [ ! -f "$ENV_FILE" ]; then
    error ".env fayli topilmadi: $ENV_FILE"
fi
success ".env fayli topildi"

# ── Server sozlash ────────────────────────────────────────────
info "Server sozlanmoqda: $SSH_USER@$SERVER_IP"

eval "$SSH_CMD" << 'REMOTE_SETUP'
set -e

# ── Paketlarni yangilash ──────────────────────────────────────
echo "📦 Paketlar yangilanmoqda..."
apt-get update -qq

# ── Git o'rnatish ─────────────────────────────────────────────
if ! command -v git &>/dev/null; then
    echo "📦 Git o'rnatilmoqda..."
    apt-get install -y git -qq
fi
echo "✅ Git: $(git --version)"

# ── Go o'rnatish ──────────────────────────────────────────────
if ! command -v go &>/dev/null; then
    echo "📦 Go o'rnatilmoqda..."
    GO_VERSION="1.22.3"
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
fi
export PATH=$PATH:/usr/local/go/bin
echo "✅ Go: $(go version)"

# ── Repo klonlash yoki yangilash ──────────────────────────────
if [ -d "/opt/ai-slayd-bot/.git" ]; then
    echo "🔄 Repo yangilanmoqda..."
    cd /opt/ai-slayd-bot
    git pull origin main
else
    echo "📥 Repo klonlanmoqda..."
    mkdir -p /opt
    git clone https://github.com/Ruslan-Xusenov/ai-slayd-bot.git /opt/ai-slayd-bot
fi

# ── Build qilish ─────────────────────────────────────────────
echo "🔨 Bot build qilinmoqda..."
cd /opt/ai-slayd-bot
export PATH=$PATH:/usr/local/go/bin
go mod download
go build -o ai-slayd-bot .
echo "✅ Build muvaffaqiyatli!"

# ── tmp papkasini yaratish ────────────────────────────────────
mkdir -p /opt/ai-slayd-bot/tmp

REMOTE_SETUP

success "Server asosiy sozlamalar bajarildi"

# ── .env faylini serverga yuklash ────────────────────────────
info ".env fayli serverga yuklanmoqda..."
eval "$SCP_CMD" "$ENV_FILE" "$SSH_USER@$SERVER_IP:/opt/ai-slayd-bot/.env"
success ".env fayli yuklandi"

# ── Systemd service yaratish ──────────────────────────────────
info "Systemd service yaratilmoqda..."

eval "$SSH_CMD" << 'SYSTEMD_SETUP'
cat > /etc/systemd/system/ai-slayd-bot.service << 'SERVICE'
[Unit]
Description=AI Slayd Bot - Telegram Presentation Generator
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/ai-slayd-bot
ExecStart=/opt/ai-slayd-bot/ai-slayd-bot
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=ai-slayd-bot
EnvironmentFile=/opt/ai-slayd-bot/.env

[Install]
WantedBy=multi-user.target
SERVICE

systemctl daemon-reload
systemctl enable ai-slayd-bot
systemctl restart ai-slayd-bot
sleep 2

if systemctl is-active --quiet ai-slayd-bot; then
    echo "✅ Bot muvaffaqiyatli ishga tushdi!"
    systemctl status ai-slayd-bot --no-pager -l
else
    echo "❌ Bot ishga tushmadi. Loglar:"
    journalctl -u ai-slayd-bot -n 20 --no-pager
fi
SYSTEMD_SETUP

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}  ✅  Deploy muvaffaqiyatli yakunlandi!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "  📊 Status ko'rish:  ${YELLOW}ssh $SSH_USER@$SERVER_IP 'systemctl status ai-slayd-bot'${NC}"
echo -e "  📋 Loglar ko'rish:  ${YELLOW}ssh $SSH_USER@$SERVER_IP 'journalctl -u ai-slayd-bot -f'${NC}"
echo -e "  🔄 Qayta ishlatish: ${YELLOW}ssh $SSH_USER@$SERVER_IP 'systemctl restart ai-slayd-bot'${NC}"
echo -e "  🛑 To'xtatish:      ${YELLOW}ssh $SSH_USER@$SERVER_IP 'systemctl stop ai-slayd-bot'${NC}"
echo ""
