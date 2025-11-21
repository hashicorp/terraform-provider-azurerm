// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-04-15/afdcustomdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importCdnFrontDoorCustomDomainAssociation() pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.FrontDoorCustomDomainAssociationID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
		customDomainId := afdcustomdomains.NewCustomDomainID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AssociationName)
		resp, err := client.Get(ctx, customDomainId)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return []*pluginsdk.ResourceData{}, fmt.Errorf("%s was not found", id)
			}
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: model was nil", id)
		}

		if resp.Model.Properties == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: properties was nil", id)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
