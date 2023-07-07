// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// (@jackofallops) Note: These tests require a valid GitHub token for ARM_GITHUB_ACCESS_TOKEN. This token needs the `repo` and `workflow` permissions on the referenced repositories.

type AppServiceGitHubTokenResource struct{}

func TestAccSourceControlGitHubToken_basic(t *testing.T) {
	token := ""
	if token = os.Getenv("ARM_GITHUB_ACCESS_TOKEN"); token == "" {
		t.Skip("Skipping as `ARM_GITHUB_ACCESS_TOKEN` is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_source_control_token", "test")
	r := AppServiceGitHubTokenResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(token),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("token").IsSet(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSourceControlGitHubToken_requiresImport(t *testing.T) {
	token := ""
	if token = os.Getenv("ARM_GITHUB_ACCESS_TOKEN"); token == "" {
		t.Skip("Skipping as `ARM_GITHUB_ACCESS_TOKEN` is not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_source_control_token", "test")
	r := AppServiceGitHubTokenResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(token),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(token),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func (r AppServiceGitHubTokenResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	resp, err := client.AppService.BaseClient.GetSourceControl(ctx, "GitHub")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), err
		}
		return nil, fmt.Errorf("retrieving Source Control GitHub Token")
	}
	if resp.Token == nil || *resp.Token == "" {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

func (r AppServiceGitHubTokenResource) basic(token string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource azurerm_source_control_token test {
  type  = "GitHub"
  token = "%s"
}
`, token)
}

func (r AppServiceGitHubTokenResource) requiresImport(token string) string {
	return fmt.Sprintf(`
%s

resource azurerm_source_control_token import {
  type  = azurerm_source_control_token.test.type
  token = azurerm_source_control_token.test.token
}
`, r.basic(token))
}
