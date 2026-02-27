# GoFortify Usage Guide 🛡️

This guide provides detailed examples and advanced usage scenarios for GoFortify.

## Basic Usage

The primary command for GoFortify is `init`, which starts both the security proxy and the TUI.

```bash
./gofortify init --port 5174 --backend-url http://localhost:8080
```

- `--port` (`-p`): The port GoFortify listens on (client-side).
- `--backend-url` (`-b`): The URL of your application that GoFortify is protecting.

## Citing GoFortify

If you use GoFortify in your academic research, security reports, or technical documentation, you can easily obtain citation info:

```bash
./gofortify cite
```

This will output:
- **Plain Text Citation**: Standard citation format for general use.
- **BibTeX Entry**: For inclusion in LaTeX-based documents and bibliography managers.

## Advanced Configuration

### Rate Limiting

You can customize the rate limiting threshold to suit your application's traffic profile.

```bash
./gofortify init -p 80 -b http://127.0.0.1:3000 --ratelimit 200
```
This allows 200 requests per minute per IP address before blocking further traffic.

### Custom Log Locations

If you want to store security logs in a specific directory or use different filenames:

```bash
./gofortify init 
  --sql /var/log/gofortify/sql_attempts.json 
  --xss /var/log/gofortify/xss_attempts.json 
  --ratelimit-file /var/log/gofortify/rate_limit.json
```

## Security Rule Details

GoFortify applies deep-packet inspection for the following:

- **SQL Injection**: Detects keywords like `UNION`, `SELECT`, `INSERT`, `DROP`, and patterns like `' OR 1=1--`.
- **Cross-Site Scripting (XSS)**: Detects tags like `<script>`, `onerror=`, `onload=`, and other malicious JavaScript vectors.

### How it works

1. **Interception**: Every incoming request is paused.
2. **Analysis**: Query parameters, form-data, and JSON bodies are scanned.
3. **Action**: If a threat is detected, GoFortify returns a `403 Forbidden` response and logs the event.
4. **Forwarding**: If no threat is detected, the request is forwarded to your backend.

## Interactive TUI Tips

- **Navigation**: Use arrow keys or `j`/`k` to navigate lists.
- **Switching Views**: Use `enter` to select "All Logs" from the main menu.
- **Back to Menu**: Press `esc` from the log view to return to the main menu.
- **Clearing Logs**: In the log view, press `c` to clear the current display (this does not delete the JSON files).
- **Exit**: Press `q` or `ctrl+c` at any time to quit the application.

---

For troubleshooting or bug reports, please visit the [main repository](https://github.com/EthicalGopher/GoFortify).
