// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedidentity

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2023-01-31/managedidentities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = FederatedIdentityCredentialResource{}

type FederatedIdentityCredentialResource struct{}

func (r FederatedIdentityCredentialResource) ModelObject() interface{} {
	return &FederatedIdentityCredentialResourceSchema{}
}

type FederatedIdentityCredentialResourceSchema struct {
	Audience          []string `tfschema:"audience"`
	Issuer            string   `tfschema:"issuer"`
	Name              string   `tfschema:"name"`
	ResourceGroupName string   `tfschema:"resource_group_name"`
	ResourceName      string   `tfschema:"parent_id"`
	Subject           string   `tfschema:"subject"`
}

func (r FederatedIdentityCredentialResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managedidentities.ValidateFederatedIdentityCredentialID
}

func (r FederatedIdentityCredentialResource) ResourceType() string {
	return "azurerm_federated_identity_credential"
}

func (r FederatedIdentityCredentialResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"audience": {
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			ForceNew: false,
			Required: true,
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
		},
		"issuer": {
			ForceNew: false,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"parent_id": {
			// TODO: this wants renaming to `user_assigned_identity_id` (and `resource_group_name` removing in 4.0)
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},
		"subject": {
			ForceNew: false,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r FederatedIdentityCredentialResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r FederatedIdentityCredentialResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20230131.ManagedIdentities

			var config FederatedIdentityCredentialResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			parentId, err := commonids.ParseUserAssignedIdentityID(config.ResourceName)
			if err != nil {
				return fmt.Errorf("parsing parent resource ID: %+v", err)
			}

			locks.ByID(parentId.ID())
			defer locks.UnlockByID(parentId.ID())

			id := managedidentities.NewFederatedIdentityCredentialID(subscriptionId, config.ResourceGroupName, parentId.UserAssignedIdentityName, config.Name)
			if metadata.ResourceData.IsNewResource() {
				existing, err := client.FederatedIdentityCredentialsGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var payload managedidentities.FederatedIdentityCredential
			r.mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredential(config, &payload)

			if _, err := client.FederatedIdentityCredentialsCreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r FederatedIdentityCredentialResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20230131.ManagedIdentities
			schema := FederatedIdentityCredentialResourceSchema{}

			id, err := managedidentities.ParseFederatedIdentityCredentialID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.FederatedIdentityCredentialsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.FederatedIdentityCredentialName
				schema.ResourceGroupName = id.ResourceGroupName
				parentId := commonids.NewUserAssignedIdentityID(id.SubscriptionId, id.ResourceGroupName, id.UserAssignedIdentityName)
				schema.ResourceName = parentId.ID()
				r.mapFederatedIdentityCredentialToFederatedIdentityCredentialResourceSchema(*model, &schema)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r FederatedIdentityCredentialResource) Update() sdk.ResourceFunc {
	return r.Create()
}

func (r FederatedIdentityCredentialResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20230131.ManagedIdentities

			var config FederatedIdentityCredentialResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parentId, err := commonids.ParseUserAssignedIdentityID(config.ResourceName)
			if err != nil {
				return fmt.Errorf("parsing parent resource ID: %+v", err)
			}

			locks.ByID(parentId.ID())
			defer locks.UnlockByID(parentId.ID())

			id, err := managedidentities.ParseFederatedIdentityCredentialID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.FederatedIdentityCredentialsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredentialProperties(input FederatedIdentityCredentialResourceSchema, output *managedidentities.FederatedIdentityCredentialProperties) {
	output.Audiences = input.Audience
	output.Issuer = input.Issuer
	output.Subject = input.Subject
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialPropertiesToFederatedIdentityCredentialResourceSchema(input managedidentities.FederatedIdentityCredentialProperties, output *FederatedIdentityCredentialResourceSchema) {
	output.Audience = input.Audiences
	output.Issuer = input.Issuer
	output.Subject = input.Subject
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredential(input FederatedIdentityCredentialResourceSchema, output *managedidentities.FederatedIdentityCredential) {
	if output.Properties == nil {
		output.Properties = &managedidentities.FederatedIdentityCredentialProperties{}
	}
	r.mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredentialProperties(input, output.Properties)
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialToFederatedIdentityCredentialResourceSchema(input managedidentities.FederatedIdentityCredential, output *FederatedIdentityCredentialResourceSchema) {
	if input.Properties == nil {
		input.Properties = &managedidentities.FederatedIdentityCredentialProperties{}
	}
	r.mapFederatedIdentityCredentialPropertiesToFederatedIdentityCredentialResourceSchema(*input.Properties, output)
}
