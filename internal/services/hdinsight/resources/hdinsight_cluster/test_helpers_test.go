// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package hdinsight_cluster_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HDInsightHBaseClusterResource struct{}

func (r HDInsightHBaseClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseHDInsightClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HDInsight.Clusters.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading HBase %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r HDInsightHBaseClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    hbase = "2.4"
  }

  gateway {
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_A4_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 2
    }

    zookeeper_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHBaseClusterResource) securityProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_hdinsight_hbase_cluster" "test" {
  name                = "acctesthdihbase-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Premium"
  tls_min_version     = "1.2"

  component_version {
    hbase = "2.4"
  }

  gateway {
    username = "sshuser"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size            = "Standard_A4_V2"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size               = "Standard_A4_V2"
      username              = "sshuser"
      password              = "TerrAform123!"
      target_instance_count = 1
      subnet_id             = azurerm_subnet.test.id
      virtual_network_id    = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size            = "Standard_A4_V2"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }
  }

  security_profile {
    aadds_resource_id       = azurerm_active_directory_domain_service.test.resource_id
    domain_name             = azurerm_active_directory_domain_service.test.domain_name
    domain_username         = azuread_user.test.user_principal_name
    domain_user_password    = azuread_user.test.password
    ldaps_urls              = ["ldaps://${azurerm_active_directory_domain_service.test.domain_name}:636"]
    msi_resource_id         = azurerm_user_assigned_identity.test.id
    cluster_users_group_dns = [azuread_group.test.display_name]
  }

  depends_on = [
    azurerm_virtual_network_dns_servers.test,
  ]
}
`, hdInsightsecurityProfileCommonTemplate(data), data.RandomInteger)
}

func (HDInsightHBaseClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hdi-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestsa%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}
	`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

type HDInsightHadoopClusterResource struct{}

func (r HDInsightHadoopClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseHDInsightClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HDInsight.Clusters.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Hadoop %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r HDInsightHadoopClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Standard"
  tls_min_version     = "1.2"


  component_version {
    hadoop = "3.3"
  }

  gateway {
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_A4_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 2
    }

    zookeeper_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightHadoopClusterResource) securityProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_hdinsight_hadoop_cluster" "test" {
  name                = "acctesthdihadoop-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Premium"
  tls_min_version     = "1.2"

  component_version {
    hadoop = "3.3"
  }

  gateway {
    username = "sshuser"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size            = "Standard_A4_V2"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size               = "Standard_A4_V2"
      username              = "sshuser"
      password              = "TerrAform123!"
      target_instance_count = 1
      subnet_id             = azurerm_subnet.test.id
      virtual_network_id    = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size            = "Standard_A4_V2"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }
  }

  security_profile {
    aadds_resource_id       = azurerm_active_directory_domain_service.test.resource_id
    domain_name             = azurerm_active_directory_domain_service.test.domain_name
    domain_username         = azuread_user.test.user_principal_name
    domain_user_password    = azuread_user.test.password
    ldaps_urls              = ["ldaps://${azurerm_active_directory_domain_service.test.domain_name}:636"]
    msi_resource_id         = azurerm_user_assigned_identity.test.id
    cluster_users_group_dns = [azuread_group.test.display_name]
  }

  depends_on = [
    azurerm_virtual_network_dns_servers.test,
  ]
}
`, hdInsightsecurityProfileCommonTemplate(data), data.RandomInteger)
}

func (HDInsightHadoopClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hdi-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}
	`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

type HDInsightInteractiveQueryClusterResource struct{}

func (r HDInsightInteractiveQueryClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseHDInsightClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HDInsight.Clusters.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading HDInsight Interactive Query Cluster (%s): %+v", id.String(), err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r HDInsightInteractiveQueryClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_interactive_query_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    interactive_hive = "3.1"
  }

  gateway {
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size  = "Standard_E16_V3"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_E16_V3"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 2
    }

    zookeeper_node {
      vm_size  = "Standard_E16_V3"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightInteractiveQueryClusterResource) securityProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_hdinsight_interactive_query_cluster" "test" {
  name                = "acctesthdiiq-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Premium"
  tls_min_version     = "1.2"

  component_version {
    interactive_hive = "3.1"
  }

  gateway {
    username = "sshuser"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size            = "Standard_E16_V3"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size               = "Standard_E16_V3"
      username              = "sshuser"
      password              = "TerrAform123!"
      target_instance_count = 1
      subnet_id             = azurerm_subnet.test.id
      virtual_network_id    = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size            = "Standard_E16_V3"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }
  }

  security_profile {
    aadds_resource_id       = azurerm_active_directory_domain_service.test.resource_id
    domain_name             = azurerm_active_directory_domain_service.test.domain_name
    domain_username         = azuread_user.test.user_principal_name
    domain_user_password    = azuread_user.test.password
    ldaps_urls              = ["ldaps://${azurerm_active_directory_domain_service.test.domain_name}:636"]
    msi_resource_id         = azurerm_user_assigned_identity.test.id
    cluster_users_group_dns = [azuread_group.test.display_name]
  }

  depends_on = [
    azurerm_virtual_network_dns_servers.test,
  ]
}
`, hdInsightsecurityProfileCommonTemplate(data), data.RandomInteger)
}

func (HDInsightInteractiveQueryClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hdi-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}
	`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

type HDInsightKafkaClusterResource struct{}

func (r HDInsightKafkaClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseHDInsightClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HDInsight.Clusters.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Kafka %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r HDInsightKafkaClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_kafka_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    kafka = "3.2"
  }

  gateway {
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size                  = "Standard_A4_V2"
      username                 = "acctestusrvm"
      password                 = "AccTestvdSC4daf986!"
      target_instance_count    = 3
      number_of_disks_per_node = 2
    }

    zookeeper_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightKafkaClusterResource) restProxy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azuread" {}

resource "azuread_group" "test" {
  display_name     = "acctesthdi-%d"
  security_enabled = true
}

resource "azurerm_hdinsight_kafka_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    kafka = "3.2"
  }

  gateway {
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size                  = "Standard_A4_V2"
      username                 = "acctestusrvm"
      password                 = "AccTestvdSC4daf986!"
      target_instance_count    = 3
      number_of_disks_per_node = 2
    }

    zookeeper_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    kafka_management_node {
      vm_size  = "Standard_A4_V2"
      password = "AccTestvdSC4daf986!"
    }
  }

  rest_proxy {
    security_group_id   = azuread_group.test.object_id
    security_group_name = azuread_group.test.display_name
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r HDInsightKafkaClusterResource) securityProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_hdinsight_kafka_cluster" "test" {
  name                = "acctesthdikafka-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Premium"
  tls_min_version     = "1.2"

  component_version {
    kafka = "3.2"
  }

  gateway {
    username = "sshuser"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size            = "Standard_A4_V2"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size                  = "Standard_A4_V2"
      username                 = "sshuser"
      password                 = "TerrAform123!"
      target_instance_count    = 3
      number_of_disks_per_node = 1
      subnet_id                = azurerm_subnet.test.id
      virtual_network_id       = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size            = "Standard_A4_V2"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }
  }

  security_profile {
    aadds_resource_id       = azurerm_active_directory_domain_service.test.resource_id
    domain_name             = azurerm_active_directory_domain_service.test.domain_name
    domain_username         = azuread_user.test.user_principal_name
    domain_user_password    = azuread_user.test.password
    ldaps_urls              = ["ldaps://${azurerm_active_directory_domain_service.test.domain_name}:636"]
    msi_resource_id         = azurerm_user_assigned_identity.test.id
    cluster_users_group_dns = [azuread_group.test.display_name]
  }

  depends_on = [
    azurerm_virtual_network_dns_servers.test,
  ]
}
`, hdInsightsecurityProfileCommonTemplate(data), data.RandomInteger)
}

func (HDInsightKafkaClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hdi-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}
	`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

type HDInsightSparkClusterResource struct{}

func (r HDInsightSparkClusterResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseHDInsightClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HDInsight.Clusters.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Spark %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r HDInsightSparkClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_spark_cluster" "test" {
  name                = "acctesthdi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Standard"
  tls_min_version     = "1.2"

  component_version {
    spark = "3.3"
  }

  gateway {
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_A4_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
    }

    zookeeper_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HDInsightSparkClusterResource) securityProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_hdinsight_spark_cluster" "test" {
  name                = "acctesthdispark-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_version     = "5.1"
  tier                = "Premium"
  tls_min_version     = "1.2"

  component_version {
    spark = "3.3"
  }

  gateway {
    username = "sshuser"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_url = azurerm_storage_container.test.url
    storage_account_key   = azurerm_storage_account.test.primary_access_key
    is_default            = true
  }

  roles {
    head_node {
      vm_size            = "Standard_A4_V2"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }

    worker_node {
      vm_size               = "Standard_A4_V2"
      username              = "sshuser"
      password              = "TerrAform123!"
      target_instance_count = 1
      subnet_id             = azurerm_subnet.test.id
      virtual_network_id    = azurerm_virtual_network.test.id
    }

    zookeeper_node {
      vm_size            = "Standard_A4_V2"
      username           = "sshuser"
      password           = "TerrAform123!"
      subnet_id          = azurerm_subnet.test.id
      virtual_network_id = azurerm_virtual_network.test.id
    }
  }

  security_profile {
    aadds_resource_id       = azurerm_active_directory_domain_service.test.resource_id
    domain_name             = azurerm_active_directory_domain_service.test.domain_name
    domain_username         = azuread_user.test.user_principal_name
    domain_user_password    = azuread_user.test.password
    ldaps_urls              = ["ldaps://${azurerm_active_directory_domain_service.test.domain_name}:636"]
    msi_resource_id         = azurerm_user_assigned_identity.test.id
    cluster_users_group_dns = [azuread_group.test.display_name]
  }

  depends_on = [
    azurerm_virtual_network_dns_servers.test,
  ]
}
`, hdInsightsecurityProfileCommonTemplate(data), data.RandomInteger)
}

func (HDInsightSparkClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hdi-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}
	`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccHDInsightHBaseCluster_securityProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hbase_cluster", "test")
	r := HDInsightHBaseClusterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.securityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account",
			"security_profile.0.domain_user_password",
			"gateway.0.password"),
	})
}

func testAccHDInsightHadoopCluster_securityProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_hadoop_cluster", "test")
	r := HDInsightHadoopClusterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.securityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account",
			"security_profile.0.domain_user_password",
			"gateway.0.password"),
	})
}

func testAccHDInsightInteractiveQueryCluster_securityProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_interactive_query_cluster", "test")
	r := HDInsightInteractiveQueryClusterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.securityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account",
			"security_profile.0.domain_user_password",
			"gateway.0.password"),
	})
}

func testAccHDInsightKafkaCluster_securityProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_kafka_cluster", "test")
	r := HDInsightKafkaClusterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.securityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account",
			"security_profile.0.domain_user_password",
			"gateway.0.password"),
	})
}

func testAccHDInsightSparkCluster_securityProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_spark_cluster", "test")
	r := HDInsightSparkClusterResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.securityProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("roles.0.head_node.0.password",
			"roles.0.head_node.0.vm_size",
			"roles.0.worker_node.0.password",
			"roles.0.worker_node.0.vm_size",
			"roles.0.zookeeper_node.0.password",
			"roles.0.zookeeper_node.0.vm_size",
			"storage_account",
			"security_profile.0.domain_user_password",
			"gateway.0.password"),
	})
}
