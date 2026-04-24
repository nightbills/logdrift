# logdrift

A lightweight CLI tool for tailing and filtering structured JSON logs across multiple services simultaneously.

---

## Installation

```bash
go install github.com/yourusername/logdrift@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logdrift.git
cd logdrift
go build -o logdrift .
```

---

## Usage

Tail logs from multiple services and filter by log level or field value:

```bash
# Tail logs from two services simultaneously
logdrift tail --services api,worker --level error

# Filter by a specific JSON field
logdrift tail --services api --filter "user_id=42"

# Watch a local log file with JSON output
logdrift tail --file ./app.log --level warn --pretty
```

**Example output:**

```
[api]     2024-03-12T10:45:01Z ERROR  "message":"connection timeout" "service":"api" "user_id":99
[worker]  2024-03-12T10:45:02Z WARN   "message":"queue backlog high" "service":"worker" "depth":512
```

### Flags

| Flag | Description |
|------|-------------|
| `--services` | Comma-separated list of service names to tail |
| `--level` | Minimum log level to display (`debug`, `info`, `warn`, `error`) |
| `--filter` | Filter by a specific JSON key-value pair |
| `--pretty` | Pretty-print JSON output |
| `--file` | Path to a local log file |

---

## License

[MIT](LICENSE)