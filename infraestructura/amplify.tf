resource "aws_amplify_app" "frontend" {
  name = "customermx-frontend"

  # Deployment manual via script — no repositorio necesario
  build_spec = <<-BUILD
    version: 1
    frontend:
      phases:
        build:
          commands:
            - echo "Manual deployment"
      artifacts:
        baseDirectory: /
        files:
          - '**/*'
  BUILD

  tags = { Name = "customermx-frontend" }
}

resource "aws_amplify_branch" "main" {
  app_id      = aws_amplify_app.frontend.id
  branch_name = "main"
}

# ── Dominio personalizado en Amplify ─────────────────────────────────────────
# Si ya existe (creado vía CLI), importar antes de aplicar:
#   terraform import aws_amplify_domain_association.main dzplpcoyzkvs1/customermx.com

resource "aws_amplify_domain_association" "main" {
  app_id      = aws_amplify_app.frontend.id
  domain_name = var.domain_name

  sub_domain {
    branch_name = aws_amplify_branch.main.branch_name
    prefix      = ""      # dominio raíz: customermx.com
  }

  sub_domain {
    branch_name = aws_amplify_branch.main.branch_name
    prefix      = "www"   # www.customermx.com
  }
}

# ── Parsear los valores calculados por Amplify ────────────────────────────────
locals {
  # certificate_verification_dns_record tiene formato:
  # "_abc123.customermx.com. CNAME _xyz.acm-validations.aws."
  cert_parts       = split(" CNAME ", aws_amplify_domain_association.main.certificate_verification_dns_record)
  cert_record_name = trimspace(local.cert_parts[0])
  cert_record_val  = trimspace(local.cert_parts[1])

  # dns_record del apex tiene formato: " CNAME d2xxx.cloudfront.net"
  apex_dns_record = [
    for s in aws_amplify_domain_association.main.sub_domains : s.dns_record
    if s.prefix == ""
  ][0]
  cloudfront_domain = trimspace(split("CNAME ", local.apex_dns_record)[1])
}

# ── Route 53: verificación del certificado SSL ────────────────────────────────
resource "aws_route53_record" "amplify_cert_verification" {
  zone_id = aws_route53_zone.main.zone_id
  name    = local.cert_record_name
  type    = "CNAME"
  ttl     = 300
  records = [local.cert_record_val]
}

# ── Route 53: www.customermx.com ──────────────────────────────────────────────
resource "aws_route53_record" "amplify_www" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "www.${var.domain_name}"
  type    = "CNAME"
  ttl     = 300
  records = [local.cloudfront_domain]
}

# ── Route 53: customermx.com (apex) ──────────────────────────────────────────
# CNAME no es válido en el apex — Route 53 soporta ALIAS (equivalente a ANAME)
resource "aws_route53_record" "amplify_root" {
  zone_id = aws_route53_zone.main.zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name                   = local.cloudfront_domain
    zone_id                = "Z2FDTNDATAQYW2" # Hosted zone ID global de CloudFront (constante de AWS)
    evaluate_target_health = false
  }
}
