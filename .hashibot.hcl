queued_behavior "release_commenter" "releases" {
  repo_prefix = "terraform-provider-"

  message = <<-EOF
    This has been released in [version ${var.release_version} of the provider](${var.changelog_link}). Please see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#provider-versions) or reach out if you need any assistance upgrading. As an example:
    ```hcl
    provider "${var.project_name}" {
        version = "~> ${var.release_version}"
    }
    # ... other configuration ...
    ```
  EOF
}

behavior "pull_request_size_labeler" "size" {
    label_prefix = "size/"
    label_map = {
        "size/XS" = {
            from = 0
            to = 30
        }
        "size/S" = {
            from = 31
            to = 60
        }
        "size/M" = {
            from = 61
            to = 150
        }
        "size/L" = {
            from = 151
            to = 300
        }
        "size/XL" = {
            from = 301
            to = 1000
        }
        "size/XXL" = {
            from = 1001
            to = 0
        }
    }
}
