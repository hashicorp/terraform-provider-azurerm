## Hadoop (Basic)

This module provisions a basic Hadoop cluster on Azure HDInsight using Terraform (specifically, using the `azurerm_hdinsight_cluster` resource).

Further information about the `azurerm_hdinsight_cluster` Terraform resource can be found [on the Terraform Website](https://terraform.io/docs/providers/azurerm/r/hdinsight_cluster.html); and [product information about HDInsight can be found on the Azure website](https://docs.microsoft.com/en-us/azure/hdinsight/).

### Notes

- An HDInsight Cluster can take a while to provision (typically around 15-20m)
- This example uses the same credentials for the local user account and for the web portal - which **we do not recommend in Production** (but is acceptable in an example scenario). Please ensure you use unique credentials as required, you may also wish to consider using SSH authentication for the local user accounts instead of password authentication.

### Resources

This example module provisions the following Resources:

- [`azurerm_resource_group`](https://terraform.io/docs/providers/azurerm/r/resource_group.html) - which contains the Resources
- [`azurerm_storage_account`](https://terraform.io/docs/providers/azurerm/r/storage_account.html) - which contains the Storage Container where the data's stored
- [`azurerm_storage_container`](https://terraform.io/docs/providers/azurerm/r/storage_container.html) - which contains the data from the HDInsight Cluster.
- [`azurerm_hdinsight_cluster`](https://terraform.io/docs/providers/azurerm/r/hdinsight_cluster.html) - which provisions the HDInsight Cluster.

### Variables

- `prefix` (Required) - The prefix used for all resources in this example.
- `location` (Optional, Defaults to `West Europe`) - The Azure Region in which resources should be provisioned.
- `username` - (Optional, Defaults to `sshuser`) - The username of the administrator for this HDInsight Cluster.
- `password` - (Optional, Defaults to `Password123!`) - The password associated with the administrator for this HDInsight Cluster.
- `tags` (Optional, Defaults to `{}`) - Any tags which should be assigned to resources in this example.
