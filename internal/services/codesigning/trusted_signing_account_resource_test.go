// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package codesigning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2024-09-30-preview/codesigningaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/codesigning"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type TrustedSigningAccountResource struct{}

func (a TrustedSigningAccountResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := codesigningaccounts.ParseCodeSigningAccountID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.CodeSigning.Client.CodeSigningAccounts.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccTrustedSigningAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, codesigning.TrustedSigningAccountResource{}.ResourceType(), "test")
	r := TrustedSigningAccountResource{}
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

func TestAccTrustedSigningAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, codesigning.TrustedSigningAccountResource{}.ResourceType(), "test")
	r := TrustedSigningAccountResource{}
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTrustedSigningAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, codesigning.TrustedSigningAccountResource{}.ResourceType(), "test")
	r := TrustedSigningAccountResource{}
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

func (a TrustedSigningAccountResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_trusted_signing_account" "test" {
  name                = "acctest-%[2]s"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, a.template(data), data.RandomString, data.Locations.Primary)
}

func (a TrustedSigningAccountResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_trusted_signing_account" "test" {
  name                = "acctest-%[2]s"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Premium"
  tags = {
    key = "example"
  }
}
`, a.template(data), data.RandomString, data.Locations.Primary)
}

func (a TrustedSigningAccountResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
