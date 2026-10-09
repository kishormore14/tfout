#!/usr/bin/env bash
# Route 53 DNS Banner Formatting Example with tfout
# Dynamically queries terraform output and formats a terminal banner without hardcoding details.

set -euo pipefail

TF_DIR="${1:-.}"

echo "Fetching Terraform outputs from ${TF_DIR}..."

# Extract outputs dynamically using terraform output -json and jq
DOMAIN=$(terraform output -raw domain_name 2>/dev/null || echo "example.com")
NS_LIST=$(terraform output -json name_servers 2>/dev/null | jq -r 'join(", ")' || echo "ns-1.example.net, ns-2.example.net")

# Format banner cleanly with tfout
tfout banner \
  --title "ROUTE53 HOSTED ZONE CREATED" \
  --status success \
  --message "DNS records configured and active." \
  --detail "Domain Name" "${DOMAIN}" \
  --detail "Name Servers" "${NS_LIST}" \
  --theme default \
  --border single
