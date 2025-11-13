# Miner - Standalone Database Manager

A cross-platform system tray application that bundles Adminer database manager with FrankenPHP server. Built with Go for easy distribution across Windows, macOS, and Linux.

## Features

- **System Tray Application**: Runs in the background with easy access from system tray
- **Custom Domain Routing**: Automatically configures `miner.local` in hosts file
- **CLI Commands**: Provides `php`, `fphp`, and `miner` commands globally
- **Auto-start**: Optional automatic startup on boot
- **Cross-platform**: Builds for Windows, macOS (Intel & ARM), and Linux

## Installation

1. Build or download the Miner binary for your platform.
2. Run the install (requires administrator/root privileges):
   ```bash
   sudo miner install
   ```
   This will:
   - Auto-install FrankenPHP if missing (macOS/Linux)
   - Add `miner.local` to your hosts file pointing to `127.0.0.1`
   - Register CLI commands (`php`, `fphp`, `miner`) in your PATH
   - Install auto-start service (optional)
3. Start Miner:
   ```bash
   miner
   ```
   No privileges needed for normal operation.
4. Access Adminer at http://miner.local:88 (non-privileged port)

## Usage

### Commands

```bash
miner              # Start the Miner system tray application
miner install      # Install and configure Miner (requires admin/root)
miner uninstall    # Remove Miner configuration (requires admin/root)
miner help         # Show help message
miner version      # Show version information
```

### System Tray Menu

- **Open Adminer**: Opens http://miner.local in your default browser
- **Start/Stop Server**: Toggle the Adminer server
- **Auto-start on Boot**: Enable/disable automatic startup
- **Uninstall**: Removes all configuration (hosts entry, CLI commands, auto-start)
- **Quit**: Exit the application

### CLI Commands

After installation, three commands are available globally:

```bash
# Run PHP scripts using FrankenPHP
php script.php

# Access FrankenPHP directly
fphp --help

# Open Adminer in browser
miner
```

## Uninstallation

To completely remove Miner:

```bash
sudo miner uninstall
```

This removes:
- Hosts file entry for `miner.local`
- CLI commands from PATH
- Auto-start service

## Configuration

- **Port**: 80 (no port needed in URL)
- **Domain**: miner.local
- **Server Address**: 127.0.0.1
- **Assets**: Located in the application directory

## Building from Source

```bash
# Clone the repository
git clone https://github.com/4nkitd/miner.git
cd miner

# Install dependencies
go mod download

# Build for your platform
go build -o miner cmd/miner/main.go

# Or build for all platforms
make build-all
```

## Architecture

- **Go**: Core application and system integration
- **FrankenPHP**: External PHP server invoked via generated Caddyfile (auto-installed if missing)
- **Adminer**: Database management interface
- **systray**: Cross-platform system tray support
- **kardianos/service**: Auto-start service management

## TODO

- [ ] Validate PHP execution across platforms (Linux/macOS variants)
- [ ] Embed FrankenPHP or provide offline bundle
- [ ] Custom application icons for each platform
- [ ] Installer packages (MSI, PKG, DEB/RPM)
- [ ] Configuration file support
- [ ] Multi-database connection profiles
- [ ] Plugin system for Adminer extensions

## Requirements

- **Administrator/Root**: Only required on first run to configure hosts file and PATH
- **Port 80**: Must be available (requires elevated privileges to bind)

## License

MIT

## Credits

- [Adminer](https://www.adminer.org/) - Database management in a single PHP file
- [FrankenPHP](https://frankenphp.dev/) - Modern PHP app server in Go
- [systray](https://github.com/getlantern/systray) - Cross-platform system tray
