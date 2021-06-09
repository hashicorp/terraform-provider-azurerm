## Example: Azure HDInsight

This example is used to demonstrate how to provision a HDInsight Cluster that uses Private Link Endpoints.

This example will provision:
- A Resource Group
- A StorageV2 Storage Account
- A Storage Data Lake Gen2 Filesystem
- A Virtual Network with one subnet:
  - One subnet so the HDInsight Cluster can still access its outbound dependencies
- A Service Principal to be used as the HDInsight cluster storage identity
  - Assigns the *Storage Blob Data Owner* role to the Service Principal
- A Subnet
  - A Subnet Network Security Group
  - A Subnet Network Security Group association
- A Puplic IP
- A NAT Gateway
  - A NAT Gateway Public IP association
  - A NAT Gateway Subnet association
- HDInsight Hadoop cluster with a user managed service principal and a *private link*

## Usage

- Provide values to all variables.
- Create with `terraform apply`
- Destroy all with `terraform destroy`

## Example Usage with Private Link Enabled

Enabling Private Link requires extensive networking knowledge to setup correctly. When `private_link_enabled` is set to `true`, internal standard load balancers (SLB) are created, and an Azure Private Link Service is provisioned for each SLB. The Private Link Service is what allows you to access the HDInsight cluster from private endpoints.

Standard load balancers do not automatically provide the public outbound NAT like basic load balancers do. You must provide your own NAT solution, such as Virtual Network NAT, for outbound dependencies. Your HDInsight cluster still needs access to its outbound dependencies. If these outbound dependencies are not allowed, the cluster creation will fail. You can read more about enabling a private links in a HDInsight Cluster [here](https://docs.microsoft.com/en-us/azure/hdinsight/hdinsight-private-link#how-to-create-clusters).

**NOTE:** When using the example code in `main.tf` for a `azurerm_hdinsight_interactive_query_cluster` or a `azurerm_hdinsight_spark_cluster` you will need to update the `resource type`(e.g. `azurerm_hdinsight_hadoop_cluster` to `azurerm_hdinsight_spark_cluster`), `name`, `cluster_version` and `component_version` of the cluster. Other than those four values the configuration HCL is identical.

Example of Updating to provision a Spark Cluster instead of a Hadoop Cluster:

```hcl
resource "azurerm_hdinsight_spark_cluster" "test" {
  depends_on = [azurerm_role_assignment.test, azurerm_nat_gateway.test]

  name                = "${var.prefix}-spark"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    spark = "2.4"
  }
 ...
}
```
