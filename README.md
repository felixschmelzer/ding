# ding

Prefix any command with `ding` and get a Telegram notification when it finishes — with exit status and duration.

```
ding make build

✅ Done
make build
Exit: 0 | Duration: 2m 14s | Finished: 14:32:01
```

## Install

**Homebrew:**

```bash
brew tap felixschmelzer/tap
brew install ding
```

**curl** (no Go required):

```bash
curl -fsSL https://github.com/felixschmelzer/ding/releases/latest/download/ding-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o /tmp/ding && chmod +x /tmp/ding && sudo mv /tmp/ding /usr/local/bin/ding
```

<details>
<summary>Pick your platform manually</summary>

| Platform | Binary |
|----------|--------|
| macOS Apple Silicon | `ding-darwin-arm64` |
| macOS Intel | `ding-darwin-amd64` |
| Linux x86-64 | `ding-linux-amd64` |
| Linux ARM64 | `ding-linux-arm64` |

Download from [github.com/felixschmelzer/ding/releases](https://github.com/felixschmelzer/ding/releases), then:

```bash
chmod +x ding-* && sudo mv ding-* /usr/local/bin/ding
```

</details>

**With Go installed:**

```bash
go install github.com/felixschmelzer/ding@latest
```

## Setup

Run the interactive setup once:

```bash
ding -c
```

You'll need a Telegram bot token and your chat ID:

1. Message [@BotFather](https://t.me/BotFather) → `/newbot` → copy the token
2. Send a message to your bot, then open `https://api.telegram.org/bot<TOKEN>/getUpdates` and find `"id"` inside `"chat"`

Config is saved to `~/.config/ding/config.json`.

## Usage

```bash
ding make build
ding npm run test
ding ./deploy.sh
```

Output is passed through as normal. When the command finishes you get a Telegram message.

## Options

```
ding <command> [args...]

  -c, --config   Open the interactive setup
  -v, --version  Print version
  -h, --help     Show help
```
