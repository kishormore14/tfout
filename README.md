# `tfout` — Reusable Open-Source Terminal Formatter for HashiCorp Terraform

[![License: Apache-2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8.svg)](https://go.dev)

`tfout` is a lightweight, cross-platform CLI tool designed to make **HashiCorp Terraform** terminal output cleaner, more structured, and visually impressive. It transforms raw outputs, command notifications, and JSON payloads into beautifully styled terminal banners, borders, and status cards.

---

## Key Features

- 🎨 **Status Banners**: Print clean status banners for `success`, `warning`, `error`, and `info`.
- 📄 **JSON Schema Renderer**: Format structured JSON payloads from files or standard input (`stdin`).
- 🛠️ **Terraform Outputs Integration**: Read and format Terraform working directory outputs (`terraform output -json`) directly.
- 🔒 **Sensitive Value Redaction**: Automatically redact Terraform sensitive outputs (`[SENSITIVE REDACTED]`) with opt-in `--show-sensitive` override.
- 🛡️ **Terminal Injection Protection**: Sanitize ANSI control sequences to prevent terminal escape-sequence injection attacks.
- 🎨 **Flexible Styling**: Choose between `auto`, `always`, or `never` color modes, `default`, `minimal`, or `mono` themes, and `single`, `double`, or `none` borders.
- 🌍 **Cross-Platform**: Zero external runtime dependencies. Runs on Windows, Linux, and macOS.

---

## Important Notice on Terraform Output Interception

> [!NOTE]
> `tfout` is a standalone formatting tool. Installing `tfout` **does not** automatically replace, intercept, or modify Terraform's native output during `terraform apply` or `terraform plan`. `tfout` is called post-apply (or wrapped in deployment scripts/CI pipelines) using `tfout outputs` or `tfout banner`.

---

## Prerequisites

- **Normal Usage**: Downloaded `tfout` binary. No external dependencies required for `tfout banner` or `tfout render`. `tfout outputs` requires HashiCorp Terraform binary on your system PATH.
- **Development**: Go 1.21 or higher.

---

## Installation

### Download Release Binaries

Download the appropriate binary for your platform from the GitHub Releases page:

- **Windows**: `tfout_v0.1.0_windows_amd64.exe`
- **Linux**: `tfout_v0.1.0_linux_amd64`
- **macOS (Apple Silicon)**: `tfout_v0.1.0_darwin_arm64`
- **macOS (Intel)**: `tfout_v0.1.0_darwin_amd64`

#### Verifying Release Checksums

To verify the integrity of your download:

```bash
# Linux / macOS
sha256sum -c checksums.txt

# Windows PowerShell
Get-FileHash -Algorithm SHA256 tfout_v0.1.0_windows_amd64.exe
```

### Build from Source

```bash
git clone https://github.com/tfout/tfout.git
cd tfout
go build -o tfout ./cmd/tfout
```

---

## Usage & Command Examples

### 1. `tfout banner`

Print custom formatted banners for infrastructure deployment notifications.

```bash
tfout banner --title "ROUTE53 HOSTED ZONE CREATED" \
  --status success \
  --message "Hosted zone created successfully with active name servers." \
  --detail "Domain Name" "example.com" \
  --detail "Name Servers" "ns-101.awsdns-12.com, ns-202.awsdns-25.net"
```

#### Output Example (Mono/Ascii Mode):

```text
┌────────────────────────────────────────────────────────────────────────────┐
│ [SUCCESS]  ROUTE53 HOSTED ZONE CREATED                                     │
├────────────────────────────────────────────────────────────────────────────┤
│ Hosted zone created successfully with active name servers.                 │
├────────────────────────────────────────────────────────────────────────────┤
│ Domain Name:   example.com                                                 │
│ Name Servers:  ns-101.awsdns-12.com, ns-202.awsdns-25.net                  │
└────────────────────────────────────────────────────────────────────────────┘
```

### 2. `tfout render`

Render structured JSON files or piped streams using the versioned `tfout` JSON schema.

```bash
# Render from file
tfout render --file examples/summary.json

# Render from stdin pipe
cat examples/summary.json | tfout render
```

### 3. `tfout outputs`

Read outputs directly from a Terraform project directory (`terraform output -json`).

```bash
# Format outputs in current directory
tfout outputs

# Format outputs in specific directory
tfout outputs --dir ./examples/terraform

# Reveal sensitive outputs (displays security warning)
tfout outputs --show-sensitive
```

### 4. `tfout doctor`

Check environment diagnostics and system readiness.

```bash
tfout doctor
```

### 5. `tfout version`

Display version information.

```bash
tfout version
```

---

## Formatting & Styling Options

| Flag | Values | Default | Description |
| :--- | :--- | :--- | :--- |
| `--color` | `auto`, `always`, `never` | `auto` | Control ANSI color output. Respects `NO_COLOR` env var. |
| `--theme` | `default`, `minimal`, `mono` | `default` | Visual theme selection. `mono` uses ASCII characters. |
| `--border` | `single`, `double`, `none` | `single` | Box border drawing style. |
| `--show-sensitive` | boolean | `false` | Reveal Terraform sensitive output values. |

---

## JSON Schema Specification (`v1.0`)

The schema definition for `tfout render` is located in [`schema/v1/render-schema.json`](file:///c:/Users/MIPL0027/Downloads/teraafom/schema/v1/render-schema.json).

```json
{
  "schema_version": "1.0",
  "title": "AWS VPC DEPLOYED",
  "status": "success",
  "message": "VPC deployed with public and private subnets.",
  "theme": "default",
  "border_style": "single",
  "details": {
    "VPC ID": "vpc-0a1b2c3d4e5f6g7h8",
    "CIDR Block": "10.0.0.0/16"
  },
  "items": [
    "NAT Gateways configured in us-east-1a and us-east-1b",
    "Flow logs enabled to CloudWatch"
  ]
}
```

---

## Route 53 Name Servers Example (Dynamic & Provider-Agnostic)

Instead of hardcoding domain or name server values, `tfout` allows you to format outputs dynamically in bash or PowerShell scripts:

```bash
#!/usr/bin/env bash
DOMAIN=$(terraform output -raw domain_name)
NAME_SERVERS=$(terraform output -json name_servers | jq -r 'join(", ")')

tfout banner \
  --title "ROUTE53 HOSTED ZONE CREATED" \
  --status success \
  --detail "Domain" "${DOMAIN}" \
  --detail "Name Servers" "${NAME_SERVERS}"
```

---

## Security Controls

1. **Sensitive Data Redaction**: Outputs marked `"sensitive": true` by Terraform are hidden as `[SENSITIVE REDACTED]` unless `--show-sensitive` is passed.
2. **Control Sequence Sanitization**: Prevents ANSI escape sequence injection attacks from untrusted Terraform output values or JSON inputs.
3. **Safe Subprocess Execution**: Subprocess execution for `terraform output -json` does not construct shell strings, avoiding command injection vulnerabilities.
4. **Input Boundary Enforcement**: Enforces a 10MB maximum payload limit on JSON input rendering.

---

## Known Limitations & Troubleshooting

- `tfout outputs` requires HashiCorp Terraform binary to be installed on your system PATH.
- `tfout` formats existing state outputs; it does not generate or execute `terraform apply`.
- `NO_COLOR` environment variable automatically turns off color unless overridden with `--color=always`.

---

## License

`tfout` is licensed under the [Apache License 2.0](LICENSE).
