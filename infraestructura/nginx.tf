# Route 53 - crea la hosted zone
resource "aws_route53_zone" "main" {
  name = var.domain_name
  tags = { Name = "customermx-zone" }
}

# Registro A: api.customermx.com -> IP elástica del EC2
resource "aws_route53_record" "api" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "api.${var.domain_name}"
  type    = "A"
  ttl     = 300
  records = [aws_eip.backend.public_ip]
}

# Instalación de Nginx + HTTPS vía SSH al EC2
resource "null_resource" "nginx_setup" {
  depends_on = [aws_route53_record.api]

  triggers = {
    domain = var.domain_name
    ec2_id = aws_instance.backend.id
  }

  connection {
    type        = "ssh"
    user        = "ec2-user"
    private_key = file("${path.module}/customermx-key.pem")
    host        = aws_eip.backend.public_ip
  }

  # Sube la config de Nginx renderizada con el dominio
  provisioner "file" {
    content = templatefile("${path.module}/scripts/nginx.conf.tpl", {
      api_domain = "api.${var.domain_name}"
    })
    destination = "/tmp/nginx_api.conf"
  }

  provisioner "remote-exec" {
    inline = [
      # Instalar Nginx y certbot
      "sudo dnf install -y nginx certbot python3-certbot-nginx",
      "sudo systemctl enable nginx",
      "sudo systemctl start nginx || sudo systemctl restart nginx",

      # Colocar config y recargar
      "sudo mv /tmp/nginx_api.conf /etc/nginx/conf.d/api.${var.domain_name}.conf",
      "sudo nginx -t",
      "sudo systemctl reload nginx",

      # Esperar propagación del DNS antes de solicitar el certificado
      "echo 'Esperando 90s para propagación DNS...'",
      "sleep 90",

      # Solicitar certificado TLS y redirigir HTTP->HTTPS
      "sudo certbot --nginx -d api.${var.domain_name} --non-interactive --agree-tos -m ${var.certbot_email} --redirect",

      "echo 'HTTPS listo en https://api.${var.domain_name}'",
    ]
  }
}
