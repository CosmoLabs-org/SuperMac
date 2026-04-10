# Installation

## Quick Install (macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/CosmoLabs-org/SuperMac/master/install.sh | bash
```

Installs the `mac` binary to `/usr/local/bin/`. Detects Apple Silicon (arm64) or Intel (amd64) automatically.

## Homebrew

```bash
brew tap cosmolabs-org/tap
brew install supermac
```

## Build from Source

Requires [Go 1.26+](https://go.dev/dl/) and macOS 12.0+.

```bash
git clone https://github.com/CosmoLabs-org/SuperMac.git
cd SuperMac/supermac-go
make build          # builds ./mac binary
make test           # run all tests (92 tests)
make install        # install to ~/bin/
```

## Updating

**Installed via curl script:**
```bash
curl -fsSL https://raw.githubusercontent.com/CosmoLabs-org/SuperMac/master/install.sh | bash
```
Re-running the install script replaces the binary with the latest release.

**Installed via Homebrew:**
```bash
brew upgrade supermac
```

**Built from source:**
```bash
cd SuperMac/supermac-go
git pull
make install
```

## Shell Completions

```bash
mac completion zsh  > ~/.zfunc/_mac
mac completion bash > /etc/bash_completion.d/mac
mac completion fish > ~/.config/fish/completions/mac.fish
```

## Verify Installation

```bash
mac version         # show version + build info
mac help            # list all modules
mac system info     # test a real command
```

## Optional Dependencies

Some commands need extra tools:

| Command | Dependency | Install |
|---------|-----------|---------|
| `mac bluetooth *` | blueutil | `brew install blueutil` |
| `mac dock add/remove` | dockutil | `brew install dockutil` |
| `mac audio input/output` | SwitchAudioSource | `brew install switchaudio-osx` |

All other commands work without extra dependencies.
