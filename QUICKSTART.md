# Miner - Quick Start Guide

## Prerequisites

No manual prerequisites on macOS/Linux: the `install` subcommand will auto-install FrankenPHP if it is missing.

Optional manual install (skips auto-download):
```bash
curl https://frankenphp.dev/install.sh | sh
sudo mv frankenphp /usr/local/bin/
frankenphp version
```

## Installation

### Step 1: Build Miner

First, build the binary:

```bash
go build -o build/miner cmd/miner/main.go
```

### Step 2: Install Miner (one-time setup)

Run the install command with sudo:

```bash
sudo ./build/miner install
```

This will:
- Auto-install FrankenPHP if missing (macOS/Linux)
- Add `miner.local -> 127.0.0.1` to `/etc/hosts`
- Create wrapper scripts in `./build/bin/` for `php`, `fphp`, `miner`
- Add `./build/bin/` to your system PATH
- Install auto-start service (optional)

**Important:** After installation, restart your terminal or run:
```bash
source ~/.zshrc  # or source ~/.bashrc
```

### Step 3: Start Miner

Once installed, start Miner (no sudo needed):

```bash
./build/miner
```

Or if PATH is configured correctly and you're in any directory:

```bash
miner
```

### Step 3: Access Adminer

Open your browser to: **http://miner.local:88**

(Port 88 is used to avoid binding to privileged port 80. Browsers default to port 80 when no port is specified; there is no hosts file mechanism to hide a non-standard port.)

Adminer will execute PHP via FrankenPHP once the server starts.

## CLI Commands

**Note:** These commands are only available after running `sudo ./build/miner install`

After installation and restarting your terminal:

```bash
# Run PHP scripts using FrankenPHP
php script.php
php -v

# Access FrankenPHP directly  
fphp php-server      # Start FrankenPHP server
fphp --help          # FrankenPHP help

# Open Adminer in browser (if miner is running)
miner
```

## Troubleshooting

### "command not found: php" or "command not found: fphp"

1. Make sure you ran `sudo ./build/miner install`
2. Restart your terminal or run `source ~/.zshrc`
3. Check if the bin directory is in your PATH:
   ```bash
   echo $PATH | grep miner
   ```

### CLI commands ask for sudo

This means the wrapper scripts weren't created properly. Run:
```bash
sudo ./build/miner uninstall
sudo ./build/miner install
```

## System Tray

When Miner is running, you'll see it in your system tray with options:
- Open Adminer (browser)
- Start/Stop Server
- Auto-start on Boot (toggle)
- Uninstall
- Quit

## Uninstallation

```bash
sudo ./build/miner uninstall
```

Removes all configuration and registered commands.

## Current Status

✅ **Working:**
- Subcommand architecture (install/uninstall)
- Hosts file management
- Port 80 (no port in URL needed)
- System tray UI
- Assets directory detection
- FrankenPHP auto-download (macOS/Linux)
- CLI commands (`php`, `fphp`, `miner`)

⏳ **TODO:**
- System tray icons
- Test auto-start service on boot
- Cross-platform builds and installers

## How It Works

Miner uses FrankenPHP as an external dependency:
- Automatically installed on macOS/Linux during `miner install` if missing
- On Windows, use WSL: `curl https://frankenphp.dev/install.sh | sh`
- Miner generates a Caddyfile to configure FrankenPHP
- FrankenPHP serves Adminer PHP files with PHP execution
- CLI commands wrap FrankenPHP for easy access
