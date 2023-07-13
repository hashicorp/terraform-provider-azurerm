// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement" // nolint: staticcheck
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func identityProviderImportFunc(providerType apimanagement.IdentityProviderType) *schema.ResourceImporter {
	return pluginsdk.ImporterValidatingResourceId(func(id string) error {
		parsed, err := parse.IdentityProviderID(id)
		if err != nil {
			return err
		}

		if parsed.Name != string(providerType) {
			return fmt.Errorf("this resource only supports Identity Provider Type %q", string(providerType))
		}

		return nil
	})
}
