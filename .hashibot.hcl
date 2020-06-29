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

poll "closed_issue_locker" "locker" {
  schedule             = "0 50 14 * * *"
  closed_for           = "720h" # 30 days
  max_issues           = 500
  sleep_between_issues = "5s"

  message = <<-EOF
    I'm going to lock this issue because it has been closed for _30 days_ ⏳. This helps our maintainers find and focus on the active issues.

    If you feel this issue should be reopened, we encourage creating a new issue linking back to this one for added context. If you feel I made an error 🤖 🙉 , please reach out to my human friends 👉  hashibot-feedback@hashicorp.com. Thanks!
  EOF
}
