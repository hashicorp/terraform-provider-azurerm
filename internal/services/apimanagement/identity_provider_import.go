// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/identityprovider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func identityProviderImportFunc(providerType identityprovider.IdentityProviderType) *schema.ResourceImporter {
	return pluginsdk.ImporterValidatingResourceId(func(id string) error {
		parsed, err := identityprovider.ParseIdentityProviderID(id)
		if err != nil {
			return err
		}

		if parsed.IdentityProviderName != providerType {
			return fmt.Errorf("this resource only supports Identity Provider Type %q", string(providerType))
		}

		return nil
	})
}
