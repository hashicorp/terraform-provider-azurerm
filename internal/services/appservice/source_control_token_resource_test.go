// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
	resp, err := client.AppService.ResourceProvidersClient.GetSourceControl(ctx, resourceproviders.NewSourceControlID("GitHub"))
	if err != nil || resp.Model == nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), err
		}
		return nil, fmt.Errorf("retrieving Source Control GitHub Token")
	}
	if resp.Model.Properties == nil || pointer.From(resp.Model.Properties.Token) == "" {
		return pointer.To(false), nil
	}
	return pointer.To(true), nil
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
