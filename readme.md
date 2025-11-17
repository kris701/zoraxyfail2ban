# Zoraxy Fail2ban Plugin

Test

# Installation

```
sudo -i
apt update
apt upgrade
apt install fail2ban

rm /etc/fail2ban/jail.d/*.conf
cat <<EOF >/etc/fail2ban/filter.d/zoraxy.conf
[Definition]
# Protect against scanners and script kiddies – for Zoraxy from V. 3.2.4
# FAILREGEX: Counts errors (401|403|404|429|444), ignores requests for favicon.ico, robots.txt, /api/notes/, api/renew, apple-touch-icon
failregex = \[client:\s*<HOST>\].*(GET|POST|HEAD|PUT|DELETE|OPTIONS)\s+/(?!favicon\.ico|robots\.txt|api/notes/|api/renew|apple-touch-icon(?:-[^/]+)?(?:-precomposed)?\.png)[^\s]*\s+(401|403|404|429|444)
EOF
chmod 644 /etc/fail2ban/filter.d/zoraxy.conf


cat <<EOF >/etc/fail2ban/jail.local
# /etc/fail2ban/jail.local
[DEFAULT]
ignoreip = 127.0.0.1/8 ::1 192.168.178.0/24
bantime = 24h
findtime = 1h
maxretry = 3
backend = auto

[zoraxy]
enabled = true
filter = zoraxy
logpath = /opt/zoraxy/log/*.log
maxretry = 8
findtime = 15m
# Initial ban time
bantime = 1h
# Automatically increase ban time for repeat offenders
bantime.increment = true
bantime.factor = 24
bantime.maxtime = 720h
# Set according to your system and installed firewall
# banaction = iptables-allports
banaction = nftables-allports
EOF
chmod 644 /etc/fail2ban/jail.local

systemctl restart fail2ban

```

Then

```
mkdir -p /opt/zoraxy/plugins/zoraxyfail2ban
cd /opt/zoraxy/plugins/zoraxyfail2ban
wget https://github.com/kris701/ZoraxyFail2BanPlugin/releases/download/v0.1.0/zoraxyfail2ban
chmod +x zoraxyfail2ban
```