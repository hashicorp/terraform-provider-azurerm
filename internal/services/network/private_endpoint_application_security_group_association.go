// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type PrivateEndpointApplicationSecurityGroupAssociationResource struct {
}

var (
	_ sdk.Resource = PrivateEndpointApplicationSecurityGroupAssociationResource{}
)

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
			ValidateFunc: validate.PrivateEndpointID,
		},
		"application_security_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApplicationSecurityGroupID,
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

			privateEndpointClient := metadata.Client.Network.PrivateEndpointClient
			privateEndpointId, err := parse.PrivateEndpointID(state.PrivateEndpointId)
			if err != nil {
				return err
			}

			locks.ByName(privateEndpointId.Name, "azurerm_private_endpoint")
			defer locks.UnlockByName(privateEndpointId.Name, "azurerm_private_endpoint")

			ASGClient := metadata.Client.Network.ApplicationSecurityGroupsClient
			ASGId, err := parse.ApplicationSecurityGroupID(state.ApplicationSecurityGroupId)
			if err != nil {
				return err
			}

			locks.ByName(ASGId.Name, "azurerm_application_security_group")
			defer locks.UnlockByName(ASGId.Name, "azurerm_application_security_group")

			existingPrivateEndpoint, err := privateEndpointClient.Get(ctx, privateEndpointId.ResourceGroup, privateEndpointId.Name, "")
			if err != nil && !utils.ResponseWasNotFound(existingPrivateEndpoint.Response) {
				return fmt.Errorf("checking for the presence of existing PrivateEndpoint %q: %+v", privateEndpointId, err)
			}

			if utils.ResponseWasNotFound(existingPrivateEndpoint.Response) {
				return fmt.Errorf("PrivateEndpoint %q does not exsits", privateEndpointId)
			}

			existingASG, err := ASGClient.Get(ctx, ASGId.ResourceGroup, ASGId.Name)
			if err != nil && !utils.ResponseWasNotFound(existingASG.Response) {
				return fmt.Errorf("checking for the presence of existingPrivateEndpoint %q: %+v", ASGId, err)
			}

			if utils.ResponseWasNotFound(existingASG.Response) {
				return fmt.Errorf("ApplicationSecurityGroup %q does not exsits", ASGId)
			}

			resourceId := parse.NewPrivateEndpointApplicationSecurityGroupAssociationId(*privateEndpointId, *ASGId)

			input := existingPrivateEndpoint
			ASGList := existingPrivateEndpoint.ApplicationSecurityGroups

			// flag: application security group exists in private endpoint configuration
			ASGInPE := false

			if input.PrivateEndpointProperties != nil && input.PrivateEndpointProperties.ApplicationSecurityGroups != nil {
				for _, value := range *ASGList {
					if value.ID != nil && *value.ID == ASGId.ID() {
						ASGInPE = true
						break
					}
				}
			}

			if ASGInPE {
				return fmt.Errorf("A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information.", resourceId.ID(), "azurerm_private_endpoint_application_security_group_association")
			}

			if ASGList != nil {
				*ASGList = append(*ASGList, existingASG)
				input.ApplicationSecurityGroups = ASGList
			} else {
				input.ApplicationSecurityGroups = &[]network.ApplicationSecurityGroup{
					existingASG,
				}
			}

			future, err := privateEndpointClient.CreateOrUpdate(ctx, privateEndpointId.ResourceGroup, privateEndpointId.Name, input)

			if err != nil {
				return fmt.Errorf("creating %s: %+v", privateEndpointId, err)
			}

			if err := future.WaitForCompletionRef(ctx, privateEndpointClient.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", privateEndpointId, err)
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

			privateEndpointClient := metadata.Client.Network.PrivateEndpointClient

			privateEndpointId, err := parse.PrivateEndpointID(resourceId.PrivateEndpointId.ID())
			if err != nil {
				return err
			}

			locks.ByName(privateEndpointId.Name, "azurerm_private_endpoint")
			defer locks.UnlockByName(privateEndpointId.Name, "azurerm_private_endpoint")

			ASGClient := metadata.Client.Network.ApplicationSecurityGroupsClient

			ASGId, err := parse.ApplicationSecurityGroupID(resourceId.ApplicationSecurityGroupId.ID())
			if err != nil {
				return err
			}

			locks.ByName(ASGId.Name, "azurerm_application_security_group")
			defer locks.UnlockByName(ASGId.Name, "azurerm_application_security_group")

			existingPrivateEndpoint, err := privateEndpointClient.Get(ctx, privateEndpointId.ResourceGroup, privateEndpointId.Name, "")
			if err != nil && !utils.ResponseWasNotFound(existingPrivateEndpoint.Response) {
				return fmt.Errorf("checking for the presence of existing PrivateEndpoint %q: %+v", privateEndpointId, err)
			}

			if utils.ResponseWasNotFound(existingPrivateEndpoint.Response) {
				return fmt.Errorf("PrivateEndpoint %q does not exsits", privateEndpointId)
			}

			existingASG, err := ASGClient.Get(ctx, ASGId.ResourceGroup, ASGId.Name)
			if err != nil && !utils.ResponseWasNotFound(existingASG.Response) {
				return fmt.Errorf("checking for the presence of existingPrivateEndpoint %q: %+v", ASGId, err)
			}

			if utils.ResponseWasNotFound(existingASG.Response) {
				return fmt.Errorf("ApplicationSecurityGroup %q does not exsits", ASGId)
			}

			// flag: application security group exists in private endpoint configuration
			ASGInPE := false

			input := existingPrivateEndpoint
			if input.PrivateEndpointProperties != nil && input.PrivateEndpointProperties.ApplicationSecurityGroups != nil {
				ASGList := *input.PrivateEndpointProperties.ApplicationSecurityGroups
				for _, value := range ASGList {
					if value.ID != nil && *value.ID == ASGId.ID() {
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

			privateEndpointClient := metadata.Client.Network.PrivateEndpointClient

			privateEndpointId, err := parse.PrivateEndpointID(state.PrivateEndpointId)
			if err != nil {
				return err
			}

			locks.ByName(privateEndpointId.Name, "azurerm_private_endpoint")
			defer locks.UnlockByName(privateEndpointId.Name, "azurerm_private_endpoint")

			ASGClient := metadata.Client.Network.ApplicationSecurityGroupsClient

			ASGId, err := parse.ApplicationSecurityGroupID(state.ApplicationSecurityGroupId)
			if err != nil {
				return err
			}

			locks.ByName(ASGId.Name, "azurerm_application_security_group")
			defer locks.UnlockByName(ASGId.Name, "azurerm_application_security_group")

			existingPrivateEndpoint, err := privateEndpointClient.Get(ctx, privateEndpointId.ResourceGroup, privateEndpointId.Name, "")
			if err != nil && !utils.ResponseWasNotFound(existingPrivateEndpoint.Response) {
				return fmt.Errorf("checking for the presence of existing PrivateEndpoint %q: %+v", privateEndpointId, err)
			}

			if utils.ResponseWasNotFound(existingPrivateEndpoint.Response) {
				return fmt.Errorf("PrivateEndpoint %q does not exsits", privateEndpointId)
			}

			existingASG, err := ASGClient.Get(ctx, ASGId.ResourceGroup, ASGId.Name)
			if err != nil && !utils.ResponseWasNotFound(existingASG.Response) {
				return fmt.Errorf("checking for the presence of existingPrivateEndpoint %q: %+v", ASGId, err)
			}

			if utils.ResponseWasNotFound(existingASG.Response) {
				return fmt.Errorf("ApplicationSecurityGroup %q does not exsits", ASGId)
			}

			resourceId := parse.NewPrivateEndpointApplicationSecurityGroupAssociationId(*privateEndpointId, *ASGId)

			// flag: application security group exists in private endpoint configuration
			ASGInPE := false

			input := existingPrivateEndpoint
			if input.PrivateEndpointProperties != nil && input.PrivateEndpointProperties.ApplicationSecurityGroups != nil {
				ASGList := *input.PrivateEndpointProperties.ApplicationSecurityGroups
				newASGList := make([]network.ApplicationSecurityGroup, 0)
				for idx, value := range ASGList {
					if value.ID != nil && *value.ID == ASGId.ID() {
						newASGList = append(newASGList, ASGList[:idx]...)
						newASGList = append(newASGList, ASGList[idx+1:]...)
						ASGInPE = true
						break
					}
				}
				if ASGInPE {
					input.PrivateEndpointProperties.ApplicationSecurityGroups = &newASGList
				} else {
					return fmt.Errorf("deletion failed, ApplicationSecurityGroup %q does not linked with PrivateEndpoint %q", ASGId, privateEndpointId)
				}
			}

			future, err := privateEndpointClient.CreateOrUpdate(ctx, privateEndpointId.ResourceGroup, privateEndpointId.Name, input)

			if err != nil {
				return fmt.Errorf("creating %s: %+v", privateEndpointId, err)
			}

			if err := future.WaitForCompletionRef(ctx, privateEndpointClient.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", privateEndpointId, err)
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
