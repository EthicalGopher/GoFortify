# GoFortify 🛡️

**GoFortify** is a high-performance, real-time security reverse proxy and traffic inspector designed to protect your backend services from common web vulnerabilities. It acts as a shield, intercepting incoming HTTP traffic and performing deep-packet inspection before forwarding legitimate requests to your upstream application.

---

## 🌟 Key Features

- **SQL Injection (SQLi) Protection**: Automatically detects and blocks common SQL injection patterns in query parameters, form data, and JSON bodies.
- **Cross-Site Scripting (XSS) Mitigation**: Real-time detection of XSS payloads in multiple request components.
- **Intelligent Rate Limiting**: Built-in IP-based rate limiting to prevent brute-force attacks and DDoS attempts.
- **Interactive TUI (Terminal User Interface)**: A beautiful dashboard powered by `bubbletea` for real-time traffic monitoring and threat visibility.
- **Detailed JSON Logging**: Security events and blocked attempts are logged in a structured JSON format for later analysis.
- **Zero-Config Reverse Proxy**: Seamlessly integrates into your existing architecture with minimal setup.

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install) 1.25 or higher

### Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/EthicalGopher/GoFortify.git
cd GoFortify
go build -o gofortify
```

### Usage

To start the security proxy and TUI, use the `init` command:

```bash
./gofortify init --port 5174 --backend-url http://localhost:8080
```

To see how to cite GoFortify in your research or project:

```bash
./gofortify cite
```

#### Commands & Flags

| Command | Description |
|---------|-------------|
| `init`  | Starts the security proxy and TUI dashboard |
| `cite`  | Displays citation info (BibTeX, Plain Text) |

| Flag (for `init`) | Shorthand | Default | Description |
|-------------------|-----------|---------|-------------|
| `--port` | `-p` | `5174` | Local port for the proxy to listen on |
| `--backend-url` | `-b` | `http://localhost:8080` | Upstream backend server URL |
| `--ratelimit` | `-r` | `100` | Max requests allowed per minute per IP |
| `--sql` | `-s` | `vulnerabilities/sqlInjection.json` | Path to SQL injection logs |
| `--xss` | `-x` | `vulnerabilities/xss.json` | Path to XSS logs |
| `--ratelimit-file`| `-rf` | `vulnerabilities/rate_limit.json` | Path to rate limit logs |

---

## 🖥️ Terminal Interface

GoFortify features an interactive TUI built with the Charm library. It provides:
- A live feed of all proxy traffic.
- Visual alerts for blocked security threats.
- Easy navigation between logs and the main dashboard.

---

## 🛡️ Security Engine

GoFortify uses optimized regex-based inspection for:
- **Query Strings**: `?id=1' OR 1=1--`
- **JSON Bodies**: `{"username": "<script>alert(1)</script>"}`
- **Form Data**: `user_input=union select null,null...`

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Developed with ❤️ by [EthicalGopher](https://github.com/EthicalGopher)
