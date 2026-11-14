// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_policy_set_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			id, err := policysetdefinitions.ParseProviders2PolicySetDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.DeleteAtManagementGroup(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
