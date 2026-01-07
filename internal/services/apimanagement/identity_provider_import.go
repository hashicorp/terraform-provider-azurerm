// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func identityProviderImporterValidatingIdentity(providerType identityprovider.IdentityProviderType) *schema.ResourceImporter {
	return pluginsdk.ImporterValidatingIdentityThen(&identityprovider.IdentityProviderId{}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
		parsed, err := identityprovider.ParseIdentityProviderID(d.Id())
		if err != nil {
			return nil, err
		}

		if parsed.IdentityProviderName != providerType {
			return nil, fmt.Errorf("this resource only supports Identity Provider Type %q", string(providerType))
		}

		return []*pluginsdk.ResourceData{d}, nil
	})
}
