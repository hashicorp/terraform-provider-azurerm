variable "prefix" {
  description = "The prefix used for all resources in this example"
}

variable "location" {
  description = "The Azure location where all resources in this example should be created"
}

variable "github_token" {
  description = "A personal access token with the `admin:repo hook` scope."
}

variable "repo_url" {
  description = "The repository URL. For example: https://github.com/example-user/example-repo"
}
