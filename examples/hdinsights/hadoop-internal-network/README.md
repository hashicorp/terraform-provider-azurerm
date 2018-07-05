## Hadoop (on an Internal network)

This module provisions a basic Hadoop cluster attached to an Internal Network on Azure HDInsight using Terraform (specifically, using the `azurerm_hdinsight_cluster` resource).

Further information about the `azurerm_hdinsight_cluster` Terraform resource can be found [on the Terraform Website](https://terraform.io/docs/providers/azurerm/r/hdinsight_cluster.html); and [product information about HDInsight can be found on the Azure website](https://docs.microsoft.com/en-us/azure/hdinsight/).

### Resources

This example module provisions the following Resources:

- [`azurerm_resource_group`](https://terraform.io/docs/providers/azurerm/r/resource_group.html) - which contains the Resources
- [`azurerm_storage_account`](https://terraform.io/docs/providers/azurerm/r/storage_account.html) - which contains the Storage Container where the data's stored
- [`azurerm_storage_container`](https://terraform.io/docs/providers/azurerm/r/storage_container.html) - which contains the data from the HDInsight Cluster.
- [`azurerm_virtual_network`](https://terraform.io/docs/providers/azurerm/r/virtual_network.html) - which is the internal network in which the HDInsight Cluster will be provisioned.
- [`azurerm_subnet`](https://terraform.io/docs/providers/azurerm/r/subnet.html) - which is the subnet within the Virtual Network in which the HDInsight Cluster will be provisioned.
- [`azurerm_hdinsight_cluster`](https://terraform.io/docs/providers/azurerm/r/hdinsight_cluster.html) - which provisions the HDInsight Cluster.

### Variables

- `prefix` (Required) - The prefix used for all resources in this example.
- `location` (Optional, Defaults to `West Europe`) - The Azure Region in which resources should be provisioned.
- `username` - (Optional, Defaults to `sshuser`) - The username of the administrator for this HDInsight Cluster.
- `password` - (Optional, Defaults to `Password123!`) - The password associated with the administrator for this HDInsight Cluster.
- `tags` (Optional, Defaults to `{}`) - Any tags which should be assigned to resources in this example.
