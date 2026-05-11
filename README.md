# Rynx CLI

A zero-config CLI tool for Git and Jira integration.

## Features

- **Auth**: Manage Jira credentials.
- **Tickets**: List and filter Jira tickets directly from the terminal.
- **Init**: Initialize project-specific configurations.

## Installation

### macOS (Homebrew)
```bash
brew tap aryanwalia2003/homebrew-tap
brew install rynx
```

### Linux (Debian/Ubuntu)
Download the latest `.deb` from the [Releases](https://github.com/aryanwalia2003/rynx/releases) page and install:
```bash
sudo dpkg -i rynx_0.1.0_linux_amd64.deb
```

### Windows
Download the latest `.zip` or binary from the [Releases](https://github.com/aryanwalia2003/rynx/releases) page. (MSI installer coming soon).

---

> [!CAUTION]
> **Unsigned Binaries:** `rynx` is currently unsigned. 
> - **macOS:** You may need to go to `System Settings > Privacy & Security` and click "Open Anyway" after the first run.
> - **Windows:** Click "More info" > "Run anyway" on the SmartScreen prompt.

> [!WARNING]
> **Sudo Blocking:** For security and to prevent permission issues with your config files, `rynx` will block execution if run with `sudo`. Always run as a normal user.

## Development

### Prerequisites

- Go (1.21+)

### Building

```bash
make build
```

### Running

```bash
./bin/dev ping
```

## Roadmap

See [ROADMAP.md](ROADMAP.md) for future plans.
