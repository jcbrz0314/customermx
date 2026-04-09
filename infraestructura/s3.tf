resource "random_id" "bucket_suffix" {
  byte_length = 4
}

resource "aws_s3_bucket" "events" {
  bucket = "customermx-events-${random_id.bucket_suffix.hex}"
  tags   = { Name = "customermx-events" }
}

resource "aws_s3_bucket_public_access_block" "events" {
  bucket                  = aws_s3_bucket.events.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# IAM role para que EC2 acceda a S3
resource "aws_iam_role" "ec2_s3" {
  name = "customermx-ec2-s3-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy" "ec2_s3" {
  name = "customermx-s3-policy"
  role = aws_iam_role.ec2_s3.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ]
      Resource = [
        aws_s3_bucket.events.arn,
        "${aws_s3_bucket.events.arn}/*"
      ]
    }]
  })
}

resource "aws_iam_instance_profile" "ec2_s3" {
  name = "customermx-ec2-profile"
  role = aws_iam_role.ec2_s3.name
}
