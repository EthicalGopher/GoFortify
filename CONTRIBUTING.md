# Contributing to GoFortify 🛡️

Thank you for your interest in improving GoFortify! We appreciate all contributions, from bug reports to new security rules.

## How to Contribute

1. **Fork the Repository**: Create your own copy of the repository.
2. **Create a Branch**: `git checkout -b feature/your-feature-name`.
3. **Make Changes**: Implement your changes and ensure your code follows Go best practices.
4. **Test Your Changes**: Verify that the security proxy and TUI work as expected.
5. **Submit a Pull Request**: Provide a clear description of your changes and why they are necessary.

## Security Rules

If you're adding new regex patterns for SQLi or XSS, please:
- Provide examples of payloads the rule detects.
- Ensure the regex is optimized for performance (minimizing backtracking).
- Test against false positives (legitimate traffic being blocked).

## Code Style

- Use `go fmt` to format your code.
- Add descriptive documentation to all exported functions, types, and constants.
- Follow the existing project structure.

---

Questions? Feel free to open an issue or reach out to the maintainers!
