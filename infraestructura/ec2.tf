# SSH key pair generado por Terraform
resource "tls_private_key" "main" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "main" {
  key_name   = "customermx-key"
  public_key = tls_private_key.main.public_key_openssh
}

resource "local_file" "private_key" {
  content         = tls_private_key.main.private_key_pem
  filename        = "${path.module}/customermx-key.pem"
  file_permission = "0600"
}

# JWT secrets estables entre deploys
resource "random_password" "jwt_access_secret" {
  length  = 64
  special = false
}

resource "random_password" "jwt_refresh_secret" {
  length  = 64
  special = false
}

# AMI: Amazon Linux 2023 (us-west-2)
data "aws_ami" "amazon_linux_2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "backend" {
  ami                         = data.aws_ami.amazon_linux_2023.id
  instance_type               = var.instance_type
  key_name                    = aws_key_pair.main.key_name
  subnet_id                   = aws_subnet.public_1.id
  vpc_security_group_ids      = [aws_security_group.ec2.id]
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.ec2_s3.name

  user_data = <<-EOF
    #!/bin/bash
    set -e

    # Directorios de la app
    mkdir -p /opt/customermx/migrations

    # Systemd service
    cat > /etc/systemd/system/customermx.service << 'SERVICE'
[Unit]
Description=CustomerMX Backend API
After=network.target

[Service]
Type=simple
User=ec2-user
WorkingDirectory=/opt/customermx
EnvironmentFile=/opt/customermx/.env
ExecStart=/opt/customermx/customermx-api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
SERVICE

    systemctl daemon-reload
    chown -R ec2-user:ec2-user /opt/customermx
  EOF

  tags = { Name = "customermx-backend" }
}

# IP elástica para que la dirección no cambie al reiniciar
resource "aws_eip" "backend" {
  instance = aws_instance.backend.id
  domain   = "vpc"
  tags     = { Name = "customermx-eip" }
}
