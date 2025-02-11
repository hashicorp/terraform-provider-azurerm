---
subcategory: "HDInsight"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hdinsight_cluster"
description: |-
  Gets information about an existing HDInsight Cluster.

---

# Data Source: azurerm_hdinsight_cluster

Use this data source to access information about an existing HDInsight Cluster.

## Example Usage

```hcl
data "azurerm_hdinsight_cluster" "example" {
  name                = "example"
  resource_group_name = "example-resources"
}

output "https_endpoint" {
  value = data.azurerm_hdinsight_cluster.example.https_endpoint
}

output "cluster_id" {
  value = data.azurerm_hdinsight_cluster.example.cluster_id
}
```

## Argument Reference

* `name` - Specifies the name of this HDInsight Cluster.

* `resource_group_name` - Specifies the name of the Resource Group in which this HDInsight Cluster exists.

## Attributes Reference

* `id` - The fully qualified HDInsight resource ID.

* `name` - The HDInsight Cluster name.

* `cluster_id` - The HDInsight Cluster ID.

* `location` - The Azure Region in which this HDInsight Cluster exists.

* `cluster_version` - The version of HDInsights which is used on this HDInsight Cluster.

* `component_versions` - A map of versions of software used on this HDInsights Cluster.

* `gateway` - A `gateway` block as defined below.

* `edge_ssh_endpoint` - The SSH Endpoint of the Edge Node for this HDInsight Cluster, if an Edge Node exists.

* `https_endpoint` - The HTTPS Endpoint for this HDInsight Cluster.

* `kafka_rest_proxy_endpoint` - The Kafka Rest Proxy Endpoint for this HDInsight Cluster.

* `kind` - The kind of HDInsight Cluster this is, such as a Spark or Storm cluster.

* `tier` - The SKU / Tier of this HDInsight Cluster.

* `ssh_endpoint` - The SSH Endpoint for this HDInsight Cluster.

* `tls_min_version` - The minimal supported TLS version.

* `tags` - A map of tags assigned to the HDInsight Cluster.

---

A `gateway` block exports the following:

* `enabled` - Is the Ambari Portal enabled?

* `username` - The username used for the Ambari Portal.

* `password` - The password used for the Ambari Portal.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the HDInsight Cluster.
