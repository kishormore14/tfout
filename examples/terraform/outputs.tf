output "service_name" {
  description = "Name of deployed service"
  value       = local.service
}

output "environment" {
  description = "Deployment environment"
  value       = local.environment
}

output "database_password" {
  description = "Database administrative secret"
  value       = "SuperSecretPassword123!"
  sensitive   = true
}

output "endpoint_urls" {
  description = "Public API endpoints"
  value = [
    "https://api.example.com/v1",
    "https://api.example.com/v2"
  ]
}
