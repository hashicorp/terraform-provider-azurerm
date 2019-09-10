# Example: App Service Source Control Token for GitHub

This example configures App Service Source Control Token for GitHub.

> **NOTE:** Source Control Token's are configured at the subscription level, not on each App Service - as such this can only be configured Subscription-wide.

## Create a Personal Access Token

1. Sign in to GitHub, browse to [https://github.com/settings/tokens](https://github.com/settings/tokens), and then click **Generate new token**.

2. Give your token a descriptive name, select the **repo** and **admin:repo hook** scopes, and then click **Generate token**.

3. Copy the personal access token and use it to set the `github_token` variable.
