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
