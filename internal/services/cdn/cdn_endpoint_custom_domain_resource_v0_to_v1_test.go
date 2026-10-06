// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn"
)

// TestAccCdnEndpointCustomDomain_V0ToV1_530 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v5.3.0 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccCdnEndpointCustomDomain_V0ToV1_530(t *testing.T) {
	if cdn.IsCdnDeprecatedForCreation() {
		t.Skip(cdn.CreateDeprecationMessage)
	}

	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")

	r := NewCdnEndpointCustomDomainResource(os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP"), os.Getenv("ARM_TEST_DNS_ZONE"))
	r.preCheck(t)

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acceptanceRG-%[2]d/providers/Microsoft.Cdn/profiles/acceptancecdnprof%[2]d/endpoints/acceptancecdnend%[2]d/customDomains/acceptance-customdomain", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicV0(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Ensure import is a no-op so the ID stored in state remains the non-canonical one being imported
					plancheck.ExpectResourceAction(importedResourceName, plancheck.ResourceActionNoop),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acceptanceRG-%[2]d/providers/microsoft.cdn/profiles/acceptancecdnprof%[2]d/endpoints/acceptancecdnend%[2]d/customDomains/acceptance-customdomain", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.basicImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acceptanceRG-%[2]d/providers/Microsoft.Cdn/profiles/acceptancecdnprof%[2]d/endpoints/acceptancecdnend%[2]d/customDomains/acceptance-customdomain", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "5.3.0")
}

func (r CdnEndpointCustomDomainResource) basicV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_endpoint_custom_domain" "test-import" {
  name            = "acceptance-customdomain"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
}

removed {
  from = azurerm_cdn_endpoint_custom_domain.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_cdn_endpoint_custom_domain.test-import
  id = "/subscriptions/%[3]s/resourcegroups/acceptanceRG-%[2]d/providers/microsoft.cdn/profiles/acceptancecdnprof%[2]d/endpoints/acceptancecdnend%[2]d/customDomains/acceptance-customdomain"
}
`, r.template(data), data.RandomInteger, data.Subscriptions.Primary)
}

func (r CdnEndpointCustomDomainResource) basicImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cdn_endpoint_custom_domain" "test-import" {
  name            = "acceptance-customdomain"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
}
`, r.template(data))
}
