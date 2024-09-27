// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/sourcecontrol"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SourceControlResource struct {
	githubRepo
}

type githubRepo struct {
	url   string
	token string
}

func newSourceControlResource(t *testing.T) SourceControlResource {
	// - ARM_TEST_ASC_GITHUB_REPOSITORY_URL represents the user repo
	// - ARM_TEST_ASC_GITHUB_USER_TOKEN represents the github personal token with the appropriate permissions per: https://docs.microsoft.com/en-us/azure/container-registry/container-registry-tutorial-build-task#create-a-github-personal-access-token
	// Checkout https://docs.microsoft.com/en-us/azure/container-registry/container-registry-tutorial-build-task for details.
	ins := SourceControlResource{
		githubRepo: githubRepo{
			url:   os.Getenv("ARM_TEST_ASC_GITHUB_REPOSITORY_URL"),
			token: os.Getenv("ARM_TEST_ASC_GITHUB_USER_TOKEN"),
		},
	}
	if ins.url == "" || ins.token == "" {
		t.Skipf("both `ARM_TEST_ASC_GITHUB_REPOSITORY_URL` and `ARM_TEST_ASC_GITHUB_USER_TOKEN` must be set for acceptance tests!")
	}
	return ins
}

func (s SourceControlResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sourcecontrol.ParseSourceControlID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automation.SourceControl.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (s SourceControlResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (s SourceControlResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_source_control" "test" {
  name                  = "acctest-%[2]d"
  automation_account_id = azurerm_automation_account.test.id

  repository_url          = "%[3]s"
  branch                  = "main"
  folder_path             = "/runbook"
  automatic_sync          = true
  publish_runbook_enabled = true
  source_control_type     = "GitHub"
  description             = "example repo desc"

  security {
    token      = "%[4]s"
    token_type = "PersonalAccessToken"
  }
}
`, s.template(data), data.RandomInteger, s.githubRepo.url, s.githubRepo.token)
}

func (s SourceControlResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_source_control" "test" {
  name                    = "acctest-%[2]d"
  automation_account_id   = azurerm_automation_account.test.id
  repository_url          = "%[3]s"
  branch                  = "dev"
  folder_path             = "/runbook"
  automatic_sync          = true
  publish_runbook_enabled = true
  source_control_type     = "GitHub"
  description             = "example repo desc foo"

  security {
    token      = "%[4]s"
    token_type = "PersonalAccessToken"
  }
}
`, s.template(data), data.RandomInteger, s.githubRepo.url, s.githubRepo.token)
}

func TestAccSourceControl_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SourceControlResource{}.ResourceType(), "test")
	r := newSourceControlResource(t)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("branch").HasValue("main"),
			),
		},
		data.ImportStep("security"),
	})
}

func TestAccSourceControl_update(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SourceControlResource{}.ResourceType(), "test")
	r := newSourceControlResource(t)
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("branch").HasValue("main"),
			),
		},
		data.ImportStep("security"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("branch").HasValue("dev"),
			),
		},
		data.ImportStep("security"),
	})
}
