// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace_authorization_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			id, err := namespaces.ParseAuthorizationRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model RelayNamespaceAuthorizationRuleResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := namespaces.AuthorizationRule{
				Name: pointer.To(model.Name),
				Properties: &namespaces.AuthorizationRuleProperties{
					Rights: expandAuthorizationRuleRights(metadata.ResourceData),
				},
			}

			if _, err := client.CreateOrUpdateAuthorizationRule(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
