package network

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-05-01/network"
	"time"
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

			existingPrivateEndpoint, err := privateEndpointClient.Get(ctx, privateEndpointId.ResourceGroup, privateEndpointId.Name, "")
			if err != nil && !utils.ResponseWasNotFound(existingPrivateEndpoint.Response) {
				return fmt.Errorf("checking for the presence of existingPrivateEndpoint %q: %+v", privateEndpointId, err)
			}

			if !utils.ResponseWasNotFound(existingPrivateEndpoint.Response) {
				return metadata.ResourceRequiresImport("azurerm_private_endpoint", privateEndpointId)
			}

			ASGClient := metadata.Client.Network.ApplicationSecurityGroupsClient

			ASGId, err := parse.ApplicationSecurityGroupID(state.ApplicationSecurityGroupId)
			if err != nil {
				return err
			}

			existingASG, err := ASGClient.Get(ctx, ASGId.ResourceGroup, ASGId.Name)
			if err != nil && !utils.ResponseWasNotFound(existingASG.Response) {
				return fmt.Errorf("checking for the presence of existingPrivateEndpoint %q: %+v", ASGId, err)
			}

			if !utils.ResponseWasNotFound(existingASG.Response) {
				return metadata.ResourceRequiresImport("azurerm_application_security_group", ASGId)
			}

			ASGList := existingPrivateEndpoint.ApplicationSecurityGroups

			if ASGList != nil {
				*ASGList = append(*ASGList, existingASG)
			}

			input := network.PrivateEndpoint{
				PrivateEndpointProperties: &network.PrivateEndpointProperties{
					ApplicationSecurityGroups: ASGList,
				},
			}

			future, err := privateEndpointClient.CreateOrUpdate(ctx, privateEndpointId.ResourceGroup, privateEndpointId.Name, input)

			if err != nil {
				return fmt.Errorf("creating %s: %+v", privateEndpointId, err)
			}

			if err := future.WaitForCompletionRef(ctx, privateEndpointClient.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", privateEndpointId, err)
			}

			metadata.SetID(privateEndpointId)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) Read() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) Delete() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}

func (p PrivateEndpointApplicationSecurityGroupAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	//TODO implement me
	panic("implement me")
}
