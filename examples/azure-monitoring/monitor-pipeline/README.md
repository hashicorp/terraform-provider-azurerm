# Azure Monitor Pipeline

This example provisions an AKS cluster, connects it to Azure Arc, installs the extensions required by Azure Monitor pipelines, creates a Custom Location, configures the Kubernetes prerequisites, and deploys an Azure Monitor Pipeline with the complete supported configuration.

## Usage

```shell
terraform init
terraform apply
```

Destroy the example when it is no longer needed:

```shell
terraform destroy
```