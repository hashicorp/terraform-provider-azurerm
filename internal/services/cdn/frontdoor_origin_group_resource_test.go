package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontdoorOriginGroupResource struct{}

func TestAccFrontdoorOriginGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_origin_group", "test")
	r := FrontdoorOriginGroupResource{}
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

func TestAccFrontdoorOriginGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_origin_group", "test")
	r := FrontdoorOriginGroupResource{}
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

func TestAccFrontdoorOriginGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_origin_group", "test")
	r := FrontdoorOriginGroupResource{}
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

func TestAccFrontdoorOriginGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_origin_group", "test")
	r := FrontdoorOriginGroupResource{}
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

func (r FrontdoorOriginGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := afdorigingroups.ParseOriginGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorOriginGroupsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r FrontdoorOriginGroupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-afdx-%d"
  location = "%s"
}
resource "azurerm_frontdoor_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r FrontdoorOriginGroupResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_origin_group" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id

  health_probe_settings {
    probe_interval_in_seconds = 0
    probe_path                = ""
    probe_protocol            = ""
    probe_request_type        = ""
  }

  load_balancing_settings {
    additional_latency_in_milliseconds = 0
    sample_size                        = 0
    successful_samples_required        = 0
  }

  response_based_afd_origin_error_detection_settings {
    http_error_ranges {
      begin = 0
      end   = 0
    }

    response_based_detected_error_types          = ""
    response_based_failover_threshold_percentage = 0
  }

  session_affinity_state                                         = ""
  traffic_restoration_time_to_healed_or_new_endpoints_in_minutes = 0
}
`, template, data.RandomInteger)
}

func (r FrontdoorOriginGroupResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_origin_group" "import" {
  name                 = azurerm_frontdoor_origin_group.test.name
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id

  health_probe_settings {
    probe_interval_in_seconds = 0
    probe_path                = ""
    probe_protocol            = ""
    probe_request_type        = ""
  }

  load_balancing_settings {
    additional_latency_in_milliseconds = 0
    sample_size                        = 0
    successful_samples_required        = 0
  }

  response_based_afd_origin_error_detection_settings {
    http_error_ranges {
      begin = 0
      end   = 0
    }

    response_based_detected_error_types          = ""
    response_based_failover_threshold_percentage = 0
  }

  session_affinity_state                                         = ""
  traffic_restoration_time_to_healed_or_new_endpoints_in_minutes = 0
}
`, config)
}

func (r FrontdoorOriginGroupResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_origin_group" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id

  health_probe_settings {
    probe_interval_in_seconds = 0
    probe_path                = ""
    probe_protocol            = ""
    probe_request_type        = ""
  }

  load_balancing_settings {
    additional_latency_in_milliseconds = 0
    sample_size                        = 0
    successful_samples_required        = 0
  }

  response_based_afd_origin_error_detection_settings {
    http_error_ranges {
      begin = 0
      end   = 0
    }

    response_based_detected_error_types          = ""
    response_based_failover_threshold_percentage = 0
  }

  session_affinity_state                                         = ""
  traffic_restoration_time_to_healed_or_new_endpoints_in_minutes = 0
}
`, template, data.RandomInteger)
}

func (r FrontdoorOriginGroupResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_origin_group" "test" {
  name                 = "acctest-c-%d"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id

  health_probe_settings {
    probe_interval_in_seconds = 0
    probe_path                = ""
    probe_protocol            = ""
    probe_request_type        = ""
  }

  load_balancing_settings {
    additional_latency_in_milliseconds = 0
    sample_size                        = 0
    successful_samples_required        = 0
  }

  response_based_afd_origin_error_detection_settings {
    http_error_ranges {
      begin = 0
      end   = 0
    }

    response_based_detected_error_types          = ""
    response_based_failover_threshold_percentage = 0
  }

  session_affinity_state                                         = ""
  traffic_restoration_time_to_healed_or_new_endpoints_in_minutes = 0
}
`, template, data.RandomInteger)
}
