// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managedidentity

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/federatedidentitycredentials"
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
	return federatedidentitycredentials.ValidateFederatedIdentityCredentialID
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
			client := metadata.Client.ManagedIdentity.V20241130.FederatedIdentityCredentials

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

			id := federatedidentitycredentials.NewFederatedIdentityCredentialID(subscriptionId, config.ResourceGroupName, parentId.UserAssignedIdentityName, config.Name)
			if metadata.ResourceData.IsNewResource() {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var payload federatedidentitycredentials.FederatedIdentityCredential
			r.mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredential(config, &payload)

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
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
			client := metadata.Client.ManagedIdentity.V20241130.FederatedIdentityCredentials
			schema := FederatedIdentityCredentialResourceSchema{}

			id, err := federatedidentitycredentials.ParseFederatedIdentityCredentialID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
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
			client := metadata.Client.ManagedIdentity.V20241130.FederatedIdentityCredentials

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

			id, err := federatedidentitycredentials.ParseFederatedIdentityCredentialID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredentialProperties(input FederatedIdentityCredentialResourceSchema, output *federatedidentitycredentials.FederatedIdentityCredentialProperties) {
	output.Audiences = input.Audience
	output.Issuer = input.Issuer
	output.Subject = input.Subject
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialPropertiesToFederatedIdentityCredentialResourceSchema(input federatedidentitycredentials.FederatedIdentityCredentialProperties, output *FederatedIdentityCredentialResourceSchema) {
	output.Audience = input.Audiences
	output.Issuer = input.Issuer
	output.Subject = input.Subject
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredential(input FederatedIdentityCredentialResourceSchema, output *federatedidentitycredentials.FederatedIdentityCredential) {
	if output.Properties == nil {
		output.Properties = &federatedidentitycredentials.FederatedIdentityCredentialProperties{}
	}
	r.mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredentialProperties(input, output.Properties)
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialToFederatedIdentityCredentialResourceSchema(input federatedidentitycredentials.FederatedIdentityCredential, output *FederatedIdentityCredentialResourceSchema) {
	if input.Properties == nil {
		input.Properties = &federatedidentitycredentials.FederatedIdentityCredentialProperties{}
	}
	r.mapFederatedIdentityCredentialPropertiesToFederatedIdentityCredentialResourceSchema(*input.Properties, output)
}
