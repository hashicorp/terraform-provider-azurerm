// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_endpoint_application_security_group_association

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
)

func (p Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceId, err := parse.PrivateEndpointApplicationSecurityGroupAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			privateEndpointClient := metadata.Client.Network.PrivateEndpoints

			privateEndpointId, err := privateendpoints.ParsePrivateEndpointID(resourceId.PrivateEndpointId.ID())
			if err != nil {
				return err
			}

			locks.ByName(privateEndpointId.PrivateEndpointName, "azurerm_private_endpoint")
			defer locks.UnlockByName(privateEndpointId.PrivateEndpointName, "azurerm_private_endpoint")

			ASGClient := metadata.Client.Network.ApplicationSecurityGroups

			ASGId, err := applicationsecuritygroups.ParseApplicationSecurityGroupID(resourceId.ApplicationSecurityGroupId.ID())
			if err != nil {
				return err
			}

			locks.ByName(ASGId.ApplicationSecurityGroupName, "azurerm_application_security_group")
			defer locks.UnlockByName(ASGId.ApplicationSecurityGroupName, "azurerm_application_security_group")

			existingPrivateEndpoint, err := privateEndpointClient.Get(ctx, *privateEndpointId, privateendpoints.DefaultGetOperationOptions())
			if err != nil && !response.WasNotFound(existingPrivateEndpoint.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing PrivateEndpoint %q: %+v", privateEndpointId, err)
			}

			if response.WasNotFound(existingPrivateEndpoint.HttpResponse) {
				log.Printf("%s was not found - removing from state!", *privateEndpointId)
				return metadata.MarkAsGone(resourceId)
			}

			existingASG, err := ASGClient.Get(ctx, *ASGId)
			if err != nil && !response.WasNotFound(existingASG.HttpResponse) {
				return fmt.Errorf("checking for the presence of existingPrivateEndpoint %q: %+v", ASGId, err)
			}

			if response.WasNotFound(existingASG.HttpResponse) {
				log.Printf("%s was not found - removing from state!", *ASGId)
				return metadata.MarkAsGone(resourceId)
			}

			// flag: application security group exists in private endpoint configuration
			ASGInPE := false

			input := existingPrivateEndpoint
			if input.Model != nil && input.Model.Properties != nil && input.Model.Properties.ApplicationSecurityGroups != nil {
				ASGList := *input.Model.Properties.ApplicationSecurityGroups
				for _, value := range ASGList {
					if value.Id != nil && *value.Id == ASGId.ID() {
						ASGInPE = true
						break
					}
				}
			}
			if !ASGInPE {
				log.Printf("ApplicationSecurityGroup %q does not exist in %q, removing from state.", ASGId, privateEndpointId)
				return metadata.MarkAsGone(resourceId)
			}

			state := PrivateEndpointApplicationSecurityGroupAssociationModel{
				ApplicationSecurityGroupId: ASGId.ID(),
				PrivateEndpointId:          privateEndpointId.ID(),
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}
