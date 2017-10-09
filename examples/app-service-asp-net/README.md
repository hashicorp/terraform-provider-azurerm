# ASP.NET App-Service Sample

Sample to deploy an ASP.NET application into Azure App-Services.

## Creates

1. A Resource Group
2. An App Service Plan
3. An App Service for usage with .NET
4. Deploy a simple app into (3).

## Usage

- Provide values to all variables (credentials and names).
- Create with `terraform apply`
- Destroy all with `terraform destroy --force`


## Prerequisites

- Uses `curl` for application publication. Preinstall curl in your system and add it to the system path.

**Note:** The Sample uses a local provisioner prepared for a DOS shell (trivial changes needed for a bash shell).