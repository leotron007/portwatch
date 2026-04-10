# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes with configurable rules.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a config file:

```bash
portwatch --config portwatch.yaml
```

Example `portwatch.yaml`:

```yaml
interval: 30s
alert:
  - type: stdout
rules:
  allow:
    - port: 22
    - port: 80
    - port: 443
  deny:
    - port: 8080
```

Run a one-time scan without the daemon:

```bash
portwatch scan
```

Watch for changes and log alerts to a file:

```bash
portwatch --config portwatch.yaml --log /var/log/portwatch.log
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `portwatch.yaml` | Path to config file |
| `--interval` | `30s` | Poll interval |
| `--log` | stdout | Log output destination |

## License

MIT © 2024 [yourusername](https://github.com/yourusername)