// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_endpoint_application_security_group_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
)

func (p Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state PrivateEndpointApplicationSecurityGroupAssociationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			privateEndpointClient := metadata.Client.Network.PrivateEndpoints
			privateEndpointId, err := privateendpoints.ParsePrivateEndpointID(state.PrivateEndpointId)
			if err != nil {
				return err
			}

			locks.ByName(privateEndpointId.PrivateEndpointName, "azurerm_private_endpoint")
			defer locks.UnlockByName(privateEndpointId.PrivateEndpointName, "azurerm_private_endpoint")

			ASGClient := metadata.Client.Network.ApplicationSecurityGroups
			ASGId, err := applicationsecuritygroups.ParseApplicationSecurityGroupID(state.ApplicationSecurityGroupId)
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
				return fmt.Errorf("PrivateEndpoint %q does not exist", privateEndpointId)
			}

			if existingPrivateEndpoint.Model == nil || existingPrivateEndpoint.Model.Properties == nil {
				return fmt.Errorf("model/properties for %s was nil", privateEndpointId)
			}

			existingASG, err := ASGClient.Get(ctx, *ASGId)
			if err != nil && !response.WasNotFound(existingASG.HttpResponse) {
				return fmt.Errorf("checking for the presence of existingPrivateEndpoint %q: %+v", ASGId, err)
			}

			if response.WasNotFound(existingASG.HttpResponse) {
				return fmt.Errorf("ApplicationSecurityGroup %q does not exist", ASGId)
			}

			if existingASG.Model == nil || existingASG.Model.Properties == nil {
				return fmt.Errorf("model/properties for %s was nil", ASGId)
			}

			resourceId := parse.NewPrivateEndpointApplicationSecurityGroupAssociationId(*privateEndpointId, *ASGId)

			input := existingPrivateEndpoint

			ASGList := existingPrivateEndpoint.Model.Properties.ApplicationSecurityGroups

			// flag: application security group exists in private endpoint configuration
			ASGInPE := false

			if input.Model.Properties.ApplicationSecurityGroups != nil {
				for _, value := range *ASGList {
					if value.Id != nil && *value.Id == ASGId.ID() {
						ASGInPE = true
						break
					}
				}
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if ASGInPE {
					return metadata.ResourceRequiresImport(p.ResourceType(), resourceId)
				}
			}

			if ASGList != nil {
				*ASGList = append(*ASGList, privateendpoints.ApplicationSecurityGroup{
					Id: existingASG.Model.Id,
				})
				input.Model.Properties.ApplicationSecurityGroups = ASGList
			} else {
				input.Model.Properties.ApplicationSecurityGroups = &[]privateendpoints.ApplicationSecurityGroup{
					{Id: existingASG.Model.Id},
				}
			}

			// TODO: implement callback, requires migrating to an ID implementing `resourceids.ResourceId`
			if err = privateEndpointClient.CreateOrUpdateThenPoll(ctx, *privateEndpointId, *input.Model); err != nil {
				return fmt.Errorf("creating %s: %+v", privateEndpointId, err)
			}

			metadata.SetID(resourceId)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
