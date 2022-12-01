package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-04-01-preview/service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkServiceResource struct{}

func TestAccMobileNetworkService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_service", "test")
	r := MobileNetworkServiceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkService_withQosPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_service", "test")
	r := MobileNetworkServiceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withQosPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkService_withServiceQosPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_service", "test")
	r := MobileNetworkServiceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withServiceQosPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_service", "test")
	r := MobileNetworkServiceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMobileNetworkService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_service", "test")
	r := MobileNetworkServiceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_service", "test")
	r := MobileNetworkServiceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MobileNetworkServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := service.ParseServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MobileNetwork.ServiceClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MobileNetworkServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_mobile_network" "test" {
  name                = "acctest-mn-%d"
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  mobile_country_code = "001"
  mobile_network_code = "01"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r MobileNetworkServiceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_service" "test" {
  name               = "acctest-mns-%d"
  mobile_network_id  = azurerm_mobile_network.test.id
  location           = "%s"
  service_precedence = 0

  pcc_rules {
    rule_name               = "default-rule"
    rule_precedence         = 1
    traffic_control_enabled = true

    service_data_flow_templates {
      direction      = "Uplink"
      template_name  = "IP-to-server"
      ports          = []
      protocol       = ["ip"]
      remote_ip_list = ["10.3.4.0/24"]
    }

  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkServiceResource) withQosPolicy(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_service" "test" {
  name               = "acctest-mns-%d"
  mobile_network_id  = azurerm_mobile_network.test.id
  location           = "%s"
  service_precedence = 0

  pcc_rules {
    rule_name               = "default-rule"
    rule_precedence         = 1
    traffic_control_enabled = true

    service_data_flow_templates {
      direction      = "Uplink"
      template_name  = "IP-to-server"
      ports          = []
      protocol       = ["ip"]
      remote_ip_list = ["10.3.4.0/24"]
    }

    rule_qos_policy {
      allocation_and_retention_priority_level = 9
      qos_indicator                           = 9
      preemption_capability                   = "NotPreempt"
      preemption_vulnerability                = "Preemptable"

      maximum_bit_rate {
        downlink = "1 Gbps"
        uplink   = "500 Mbps"
      }
    }
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkServiceResource) withServiceQosPolicy(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_mobile_network_service" "test" {
  name               = "acctest-mns-%d"
  mobile_network_id  = azurerm_mobile_network.test.id
  location           = "%s"
  service_precedence = 0

  pcc_rules {
    rule_name               = "default-rule"
    rule_precedence         = 1
    traffic_control_enabled = true

    service_data_flow_templates {
      direction      = "Uplink"
      template_name  = "IP-to-server"
      ports          = []
      protocol       = ["ip"]
      remote_ip_list = ["10.3.4.0/24"]
    }
  }

  service_qos_policy {
    allocation_and_retention_priority_level = 9
    qos_indicator                           = 9
    preemption_capability                   = "NotPreempt"
    preemption_vulnerability                = "Preemptable"
    maximum_bit_rate {
      downlink = "1 Gbps"
      uplink   = "100 Mbps"
    }
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkServiceResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_service" "import" {
  name               = azurerm_mobile_network_service.test.name
  mobile_network_id  = azurerm_mobile_network.test.id
  location           = "%s"
  service_precedence = 0

  pcc_rules {
    rule_name               = "default-rule"
    rule_precedence         = 1
    traffic_control_enabled = true

    service_data_flow_templates {
      direction      = "Uplink"
      template_name  = "IP-to-server"
      ports          = []
      protocol       = ["ip"]
      remote_ip_list = ["10.3.4.0/24"]
    }
  }
}
`, config, data.Locations.Primary)
}

func (r MobileNetworkServiceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_service" "test" {
  name               = "acctest-mns-%d"
  mobile_network_id  = azurerm_mobile_network.test.id
  location           = "%s"
  service_precedence = 0
  pcc_rules {
    rule_name               = "default-rule"
    rule_precedence         = 1
    traffic_control_enabled = true
    rule_qos_policy {
      allocation_and_retention_priority_level = 9
      qos_indicator                           = 9
      preemption_capability                   = "NotPreempt"
      preemption_vulnerability                = "Preemptable"
      guaranteed_bit_rate {
        downlink = "100 Mbps"
        uplink   = "10 Mbps"
      }
      maximum_bit_rate {
        downlink = "1 Gbps"
        uplink   = "100 Mbps"
      }
    }

    service_data_flow_templates {
      direction      = "Uplink"
      template_name  = "IP-to-server"
      ports          = []
      protocol       = ["ip"]
      remote_ip_list = ["10.3.4.0/24"]
    }
  }
  service_qos_policy {
    allocation_and_retention_priority_level = 9
    qos_indicator                           = 9
    preemption_capability                   = "NotPreempt"
    preemption_vulnerability                = "Preemptable"
    maximum_bit_rate {
      downlink = "1 Gbps"
      uplink   = "100 Mbps"
    }
  }
  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkServiceResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_mobile_network_service" "test" {
  name               = "acctest-mns-%d"
  mobile_network_id  = azurerm_mobile_network.test.id
  location           = "%s"
  service_precedence = 0
  pcc_rules {
    rule_name               = "default-rule-2"
    rule_precedence         = 1
    traffic_control_enabled = false
    rule_qos_policy {
      allocation_and_retention_priority_level = 9
      qos_indicator                           = 9
      preemption_capability                   = "MayPreempt"
      preemption_vulnerability                = "NotPreemptable"
      guaranteed_bit_rate {
        downlink = "200 Mbps"
        uplink   = "20 Mbps"
      }
      maximum_bit_rate {
        downlink = "2 Gbps"
        uplink   = "200 Mbps"
      }
    }

    service_data_flow_templates {
      direction      = "Uplink"
      template_name  = "IP-to-server"
      ports          = []
      protocol       = ["ip"]
      remote_ip_list = ["10.3.4.0/24"]
    }
  }

  service_qos_policy {
    allocation_and_retention_priority_level = 9
    qos_indicator                           = 9
    preemption_capability                   = "NotPreempt"
    preemption_vulnerability                = "Preemptable"
    maximum_bit_rate {
      downlink = "2 Gbps"
      uplink   = "200 Mbps"
    }
  }
  tags = {
    key = "update"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}
