terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.78.0"
    }
    databricks = {
      source = "databrickslabs/databricks"
      version = "0.3.7"
    }
  }
}

module "azure_databrick_workspace" {
  source = "./dataBricksWorkspace"

  rg_name = "xiaxintestRG-dataBrick"
  rg_location = "north central us"
  workspace_name = "xiaxintestDBW"
}

module "databrick_cluster" {
  source = "./dataBricksCluster"
  databrick_cluster_depends_on = module.azure_databrick_workspace.databricks_workspace_id
  databrick_workspace_id = module.azure_databrick_workspace.databricks_workspace_id
  databrick_workspace_URL = module.azure_databrick_workspace.databricks_workspace_URL
}