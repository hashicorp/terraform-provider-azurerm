---
subcategory: "HDInsight"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hdinsight_cluster"
sidebar_current: "docs-azurerm-datasource-hdinsight-cluster"
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
  value = "${data.azurerm_hdinsight_cluster.example.https_endpoint}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of this HDInsight Cluster.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which this HDInsight Cluster exists.

## Attributes Reference

* `location` - The Azure Region in which this HDInsight Cluster exists.

* `cluster_version` - The version of HDInsights which is used on this HDInsight Cluster.

* `component_versions` - A map of versions of software used on this HDInsights Cluster.

* `gateway` - A `gateway` block as defined below.

* `edge_ssh_endpoint` - The SSH Endpoint of the Edge Node for this HDInsight Cluster, if an Edge Node exists.

* `https_endpoint` - The HTTPS Endpoint for this HDInsight Cluster.

* `kind` - The kind of HDInsight Cluster this is, such as a Spark or Storm cluster.

* `tier` - The SKU / Tier of this HDInsight Cluster.

* `ssh_endpoint` - The SSH Endpoint for this HDInsight Cluster.

* `tags` - A map of tags assigned to the HDInsight Cluster.

---

A `gateway` block exports the following:

* `enabled` - Is the Ambari Portal enabled?

* `username` - The username used for the Ambari Portal.

* `password` - The password used for the Ambari Portal.
