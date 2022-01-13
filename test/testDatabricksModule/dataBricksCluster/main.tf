terraform {
  required_providers {
    databricks = {
      source = "databrickslabs/databricks"
      version = "0.3.7"
    }
  }
}

provider "databricks" {
  azure_workspace_resource_id = var.databrick_workspace_id
  host = var.databrick_workspace_URL
  //azure_workspace_resource_id = module.azure_databrick_workspace.databricks_workspace_id
  azure_client_id     = "0cb7d7f2-0f76-49ba-b2bb-c540c4f3eb37"
  azure_client_secret = "q.dK_1AndHx._x9oeiE5SHv6xlHWts23zk"
  azure_tenant_id     = "72f988bf-86f1-41af-91ab-2d7cd011db47"
}

data "databricks_node_type" "smallest" {
  local_disk = true
  depends_on = [var.databrick_cluster_depends_on]
}

data "databricks_spark_version" "latest" {
  long_term_support = true
  depends_on = [var.databrick_cluster_depends_on]
}

resource "databricks_cluster" "interactive_cluster" {
  cluster_name            = "xiaxintest_cluster"
  spark_version           = data.databricks_spark_version.latest.id
  node_type_id            = data.databricks_node_type.smallest.id
  autotermination_minutes = 120
  autoscale {
    min_workers = 1
    max_workers = 20
  }

  library {
    pypi {
      package="kedro"
    }
  }
  library {
    pypi {
      package="plotly"
    }
  }
  library {
    pypi {
      package="azure-identity==1.5.0"
    }
  }
  library {
    pypi {
      package="azure-storage-blob==12.6.0"
    }
  }
}