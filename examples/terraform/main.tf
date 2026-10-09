# Minimal example Terraform configuration (harmless, no cloud provider credentials required)
terraform {
  required_version = ">= 1.0.0"
}

locals {
  environment = "production"
  service     = "api-gateway"
}
