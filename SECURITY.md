# Security Policy

## Reporting Security Vulnerabilities

The `tfout` maintainers take security seriously. If you discover a security vulnerability in `tfout`, please report it responsibly rather than opening a public issue.

### How to Report

Please email security reports to security@example.com (or open a Private Vulnerability Report on GitHub) with:

- Description of the vulnerability and its potential impact.
- Step-by-step instructions or proof-of-concept to reproduce.
- Any suggested remediations or mitigations.

We will acknowledge receipt within 48 hours and provide regular updates until the issue is resolved and a release update is published.

## Security Architecture & Design Principles

1. **Sensitive Data Redaction**:
   - Terraform output values marked `"sensitive": true` are redacted by default (`[SENSITIVE REDACTED]`).
   - Revealing sensitive values requires explicit user action (`--show-sensitive`) and emits a clear warning banner.

2. **Terminal Injection Protection**:
   - All input strings (titles, messages, keys, values, items) undergo control-sequence sanitization to prevent ANSI escape sequence injection attack vectors.

3. **Safe Subprocess Execution**:
   - Subprocess commands (such as `terraform output -json`) are executed directly via `exec.CommandContext` without passing parameters through shell interpreters (`cmd.exe`, `sh`, `bash`), preventing command injection risks.

4. **Input Size Limits**:
   - JSON payload rendering enforces a 10MB input limit to prevent memory exhaustion / Denial of Service.
