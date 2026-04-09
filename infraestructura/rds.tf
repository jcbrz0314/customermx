resource "aws_db_subnet_group" "main" {
  name       = "customermx-db-subnet-group"
  subnet_ids = [aws_subnet.public_1.id, aws_subnet.public_2.id]
  tags       = { Name = "customermx-db-subnet-group" }
}

# Parameter group sin SSL obligatorio
resource "aws_db_parameter_group" "postgres15" {
  name   = "customermx-pg15"
  family = "postgres15"

  parameter {
    name         = "rds.force_ssl"
    value        = "0"
    apply_method = "pending-reboot"
  }
}

resource "aws_db_instance" "postgres" {
  identifier             = "customermx-db"
  engine                 = "postgres"
  engine_version         = "15"
  instance_class         = "db.t3.micro"
  allocated_storage      = 20
  db_name                = var.db_name
  username               = var.db_username
  password               = var.db_password
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = aws_db_parameter_group.postgres15.name
  publicly_accessible    = true
  skip_final_snapshot    = true
  deletion_protection    = false
  apply_immediately      = true

  tags = { Name = "customermx-db" }
}
