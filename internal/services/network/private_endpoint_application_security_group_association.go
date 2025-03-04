// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PrivateEndpointApplicationSecurityGroupAssociationResource struct{}

var _ sdk.Resource = PrivateEndpointApplicationSecurityGroupAssociationResource{}

type PrivateEndpointApplicationSecurityGroupAssociationModel struct {
	PrivateEndpointId          string `tfschema:"private_endpoint_id"`
	ApplicationSecurityGroupId string `tfschema:"application_security_group_id"`
}

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"private_endpoint_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: privateendpoints.ValidatePrivateEndpointID,
		},
		"application_security_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
		},
	}
}

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) ModelObject() interface{} {
	return &PrivateEndpointApplicationSecurityGroupAssociationModel{}
}

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) ResourceType() string {
	return "azurerm_private_endpoint_application_security_group_association"
}

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) Create() sdk.ResourceFunc {
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
				return fmt.Errorf("PrivateEndpoint %q does not exsits", privateEndpointId)
			}

			if existingPrivateEndpoint.Model == nil || existingPrivateEndpoint.Model.Properties == nil {
				return fmt.Errorf("model/properties for %s was nil", privateEndpointId)
			}

			existingASG, err := ASGClient.Get(ctx, *ASGId)
			if err != nil && !response.WasNotFound(existingASG.HttpResponse) {
				return fmt.Errorf("checking for the presence of existingPrivateEndpoint %q: %+v", ASGId, err)
			}

			if response.WasNotFound(existingASG.HttpResponse) {
				return fmt.Errorf("ApplicationSecurityGroup %q does not exsits", ASGId)
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

			if ASGInPE {
				return fmt.Errorf("A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information.", resourceId.ID(), "azurerm_private_endpoint_application_security_group_association")
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

			if err = privateEndpointClient.CreateOrUpdateThenPoll(ctx, *privateEndpointId, *input.Model); err != nil {
				return fmt.Errorf("creating %s: %+v", privateEndpointId, err)
			}

			metadata.SetID(resourceId)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) Read() sdk.ResourceFunc {
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
				return fmt.Errorf("PrivateEndpoint %q does not exsits", privateEndpointId)
			}

			existingASG, err := ASGClient.Get(ctx, *ASGId)
			if err != nil && !response.WasNotFound(existingASG.HttpResponse) {
				return fmt.Errorf("checking for the presence of existingPrivateEndpoint %q: %+v", ASGId, err)
			}

			if response.WasNotFound(existingASG.HttpResponse) {
				return fmt.Errorf("ApplicationSecurityGroup %q does not exsits", ASGId)
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
				log.Printf("ApplicationSecurityGroup %q does not exsits in %q, removing from state.", ASGId, privateEndpointId)
				err := metadata.MarkAsGone(resourceId)
				if err != nil {
					return err
				}
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

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) Delete() sdk.ResourceFunc {
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
				return fmt.Errorf("PrivateEndpoint %q does not exsits", privateEndpointId)
			}

			if existingPrivateEndpoint.Model == nil || existingPrivateEndpoint.Model.Properties == nil {
				return fmt.Errorf("model/properties for %s was nil", privateEndpointId)
			}

			existingASG, err := ASGClient.Get(ctx, *ASGId)
			if err != nil && !response.WasNotFound(existingASG.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing %q: %+v", ASGId, err)
			}

			if response.WasNotFound(existingASG.HttpResponse) {
				return fmt.Errorf("ApplicationSecurityGroup %q does not exsits", ASGId)
			}

			resourceId := parse.NewPrivateEndpointApplicationSecurityGroupAssociationId(*privateEndpointId, *ASGId)

			// flag: application security group exists in private endpoint configuration
			ASGInPE := false

			input := existingPrivateEndpoint
			if input.Model != nil && input.Model.Properties != nil && input.Model.Properties.ApplicationSecurityGroups != nil {
				ASGList := *input.Model.Properties.ApplicationSecurityGroups
				newASGList := make([]privateendpoints.ApplicationSecurityGroup, 0)
				for idx, value := range ASGList {
					if value.Id != nil && *value.Id == ASGId.ID() {
						newASGList = append(newASGList, ASGList[:idx]...)
						newASGList = append(newASGList, ASGList[idx+1:]...)
						ASGInPE = true
						break
					}
				}
				if ASGInPE {
					input.Model.Properties.ApplicationSecurityGroups = &newASGList
				} else {
					return fmt.Errorf("deletion failed, ApplicationSecurityGroup %q is not linked with PrivateEndpoint %q", ASGId, privateEndpointId)
				}
			}

			if err = privateEndpointClient.CreateOrUpdateThenPoll(ctx, *privateEndpointId, *input.Model); err != nil {
				return fmt.Errorf("creating %s: %+v", privateEndpointId, err)
			}

			metadata.SetID(resourceId)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.PrivateEndpointApplicationSecurityGroupAssociationIDValidation
}
