output "ec2_public_ip" {
  description = "IP pública del servidor backend"
  value       = aws_eip.backend.public_ip
}

output "ec2_public_dns" {
  description = "DNS público del servidor backend"
  value       = aws_eip.backend.public_dns
}

output "backend_url" {
  description = "URL del backend (IP directa)"
  value       = "http://${aws_eip.backend.public_ip}:8080"
}

output "api_url" {
  description = "URL pública de la API con HTTPS"
  value       = "https://api.${var.domain_name}"
}

output "route53_nameservers" {
  description = "Nameservers de Route 53 — apunta tu dominio a estos en tu registrar"
  value       = aws_route53_zone.main.name_servers
}

output "rds_endpoint" {
  description = "Host del RDS PostgreSQL"
  value       = aws_db_instance.postgres.address
}

output "rds_port" {
  description = "Puerto del RDS"
  value       = aws_db_instance.postgres.port
}

output "s3_bucket" {
  description = "Nombre del bucket S3"
  value       = aws_s3_bucket.events.bucket
}

output "amplify_app_id" {
  description = "ID de la app en Amplify"
  value       = aws_amplify_app.frontend.id
}

output "amplify_url" {
  description = "URL del frontend en Amplify"
  value       = "https://main.${aws_amplify_app.frontend.id}.amplifyapp.com"
}

output "frontend_url" {
  description = "URL pública del frontend con dominio personalizado"
  value       = "https://${var.domain_name}"
}

output "ssh_command" {
  description = "Comando SSH para conectarte al servidor"
  value       = "ssh -i infraestructura/customermx-key.pem ec2-user@${aws_eip.backend.public_ip}"
}

output "jwt_access_secret" {
  description = "JWT access secret"
  value       = random_password.jwt_access_secret.result
  sensitive   = true
}

output "jwt_refresh_secret" {
  description = "JWT refresh secret"
  value       = random_password.jwt_refresh_secret.result
  sensitive   = true
}
