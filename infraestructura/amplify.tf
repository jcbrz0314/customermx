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
