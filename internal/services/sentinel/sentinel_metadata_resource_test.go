// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	sentinelmetadata "github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MetadataResource struct{}

func (r MetadataResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Sentinel.MetadataClient
	id, err := sentinelmetadata.ParseMetadataID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func TestAccMetadata_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_metadata", "test")
	r := MetadataResource{}

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

func TestAccMetadata_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_metadata", "test")
	r := MetadataResource{}

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

func TestAccMetadata_completeWithSolution(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_metadata", "test")
	r := MetadataResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeWithSolution(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMetadata_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_metadata", "test")
	r := MetadataResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MetadataResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_metadata" "test" {
  name         = "acctest"
  workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  content_id   = azurerm_sentinel_alert_rule_nrt.test.name
  kind         = "AnalyticsRule"
  parent_id    = azurerm_sentinel_alert_rule_nrt.test.id
}
`, SentinelAlertRuleNrtResource{}.basic(data))
}

func (MetadataResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_metadata" "test" {
  name                       = "acctest"
  workspace_id               = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  content_id                 = azurerm_sentinel_alert_rule_nrt.test.name
  kind                       = "AnalyticsRule"
  parent_id                  = azurerm_sentinel_alert_rule_nrt.test.id
  providers                  = ["testprovider1", "testprovider2"]
  preview_images             = ["firstImage.png"]
  preview_images_dark        = ["firstImageDark.png"]
  content_schema_version     = "2.0"
  custom_version             = "1.0"
  version                    = "1.0.0"
  threat_analysis_tactics    = ["Reconnaissance", "CommandAndControl"]
  threat_analysis_techniques = ["T1548", "t1548.001"]

  source {
    kind = "Solution"
    name = "test Solution"
    id   = "b688a130-76f4-4a07-bf57-762222a3cadf"
  }

  author {
    name  = "test user"
    email = "acc@test.com"
  }

  support {
    name  = "acc test"
    email = "acc@test.com"
    link  = "https://acc.test.com"
    tier  = "Partner"
  }

  dependency = jsonencode({
    operator = "AND",
    criteria = [
      {
        operator = "OR",
        criteria = [
          {
            contentId = "045d06d0-ee72-4794-aba4-cf5646e4c756",
            kind      = "DataConnector",
          },
          {
            contentId = "dbfcb2cc-d782-40ef-8d94-fe7af58a6f2d",
            kind      = "DataConnector"
          },
          {
            contentId = "de4dca9b-eb37-47d6-a56f-b8b06b261593",
            kind      = "DataConnector",
            version   = "2.0"
          }
        ]
      },
      {
        kind      = "Playbook",
        contentId = "31ee11cc-9989-4de8-b176-5e0ef5c4dbab",
        version   = "1.0"
      },
      {
        kind      = "Parser",
        contentId = "21ba424a-9438-4444-953a-7059539a7a1b"
      }
    ]
  })

}
`, SentinelAlertRuleNrtResource{}.basic(data))
}

func (MetadataResource) completeWithSolution(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_metadata" "test" {
  name                       = "acctest"
  workspace_id               = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  kind                       = "Solution"
  content_id                 = azurerm_sentinel_alert_rule_nrt.test.id
  parent_id                  = azurerm_sentinel_alert_rule_nrt.test.id
  providers                  = ["testprovider1", "testprovider2"]
  preview_images             = ["firstImage.png"]
  preview_images_dark        = ["firstImageDark.png"]
  content_schema_version     = "2.0"
  custom_version             = "1.0"
  threat_analysis_tactics    = ["Reconnaissance", "CommandAndControl"]
  threat_analysis_techniques = ["T1548", "t1548.001"]
  first_publish_date         = "2021-06-24"
  last_publish_date          = "2021-07-24"
  version                    = "1.0.0"
  source {
    kind = "Solution"
    name = "test Solution"
    id   = azurerm_sentinel_alert_rule_nrt.test.id
  }

  author {
    name  = "test user"
    email = "acc@test.com"
  }

  support {
    name  = "acc test"
    email = "acc@test.com"
    link  = "https://acc.test.com"
    tier  = "Partner"
  }

  category {
    domains   = ["Application"]
    verticals = ["Healthcare"]
  }

  dependency = jsonencode({
    operator = "AND",
    criteria = [
      {
        operator = "OR",
        criteria = [
          {
            contentId = "dbfcb2cc-d782-40ef-8d94-fe7af58a6f2d",
            kind      = "DataConnector"
          },
          {
            contentId = "de4dca9b-eb37-47d6-a56f-b8b06b261593",
            kind      = "DataConnector",
            version   = "2.0"
          }
        ]
      },
      {
        kind      = "Playbook",
        contentId = "31ee11cc-9989-4de8-b176-5e0ef5c4dbab",
        version   = "1.0"
      },
      {
        kind      = "Parser",
        contentId = "21ba424a-9438-4444-953a-7059539a7a1b"
      }
    ]
  })

}
`, SentinelAlertRuleNrtResource{}.basic(data))
}
