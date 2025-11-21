# Zoraxy Fail2ban Plugin

This is a simple plugin for [Zoraxy](https://github.com/tobychui/zoraxy) that enables you to use the [Fail2Ban](https://github.com/fail2ban/fail2ban) daemon to block spam bots and scapers. This plugin is mostly just the visual interface in Zoraxy, which includes some simple controls such as updating configs or manually banning IPs.

<img width="1132" height="840" alt="image" src="https://github.com/user-attachments/assets/8b4726e9-8373-4b40-a4c1-8c5604717211" />

# Installation
To use this plugin, you must first have fail2ban installed.

## Fail2Ban

All the following steps is assumed to be run with `sudo -i` or as a root user.

First install fail2ban from the APT repository:

```bash
apt update
apt upgrade
apt install fail2ban
```

Then you can create the filter that we need to parse Zoraxy log files. You can also change the config values later in Zoraxy.

```bash
rm /etc/fail2ban/jail.d/*.conf
cat <<EOF >/etc/fail2ban/filter.d/zoraxy.conf
[Definition]
# Protect against scanners and script kiddies – for Zoraxy from V. 3.2.4
# FAILREGEX: Counts errors (401|403|404|429|444), ignores requests for favicon.ico, robots.txt, /api/notes/, api/renew, apple-touch-icon
failregex = \[client:\s*<HOST>\].*(GET|POST|HEAD|PUT|DELETE|OPTIONS)\s+/(?!favicon\.ico|robots\.txt|api/notes/|api/renew|apple-touch-icon(?:-[^/]+)?(?:-precomposed)?\.png)[^\s]*\s+(401|403|404|429|444)
EOF
chmod 777 /etc/fail2ban/filter.d/zoraxy.conf
```

Then you can create the fail2ban jail config, where you can change settings to your needs. You can also change the config values later in Zoraxy.

```bash
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
chmod 777 /etc/fail2ban/jail.local
```

Finally you can restart fail2ban and allow non-root users to modify the fail2ban socket (This is needed for restarting and updating configs).

```bash
systemctl restart fail2ban

chmod 777 /var/run/fail2ban/fail2ban.sock
```

## Plugin

You can now install the Zoraxy plugin itself, by doing the following:

```bash
mkdir -p /opt/zoraxy/plugins/zoraxyfail2ban
cd /opt/zoraxy/plugins/zoraxyfail2ban
# wget <LINK_TO_LATEST_BINARY>
wget https://github.com/kris701/zoraxyfail2ban/releases/download/v1.1.1/zoraxyfail2ban
chmod +x zoraxyfail2ban
```

Then you can restart your Zoraxy server or service and you should be able to see the new plugin in the sidebar.

# Development

Execute run script `devRun.ps1`.
Then run the server with `./zoraxy -dev=true -noauth=true -port=:8564`
You can then run `devRun.ps1` whenever you want to update the binary.
The script needs WSL installed, and it launches a wsl process for the Zoraxy server.
You can use `devKill.ps1` to kill the server again.

This is a rather rudimentary dev system, if anyone can figure out to set up a propper Docker environment with Zoraxy and Fail2ban working, i would greatly appreciate the help :)
