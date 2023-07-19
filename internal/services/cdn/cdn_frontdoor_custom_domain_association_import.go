// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"

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

		client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AssociationName)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.AFDDomainProperties == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `AFDDomainProperties` was nil", id)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
