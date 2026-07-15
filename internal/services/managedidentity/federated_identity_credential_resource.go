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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/federatedidentitycredentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name federated_identity_credential -service-package-name managedidentity -properties "name" -compare-values "subscription_id:user_assigned_identity_id,resource_group_name:user_assigned_identity_id,user_assigned_identity_name:user_assigned_identity_id" -test-sequential

var (
	_ sdk.Resource             = FederatedIdentityCredentialResource{}
	_ sdk.ResourceWithIdentity = FederatedIdentityCredentialResource{}
)

type FederatedIdentityCredentialResource struct{}

func (r FederatedIdentityCredentialResource) ModelObject() interface{} {
	return &FederatedIdentityCredentialResourceSchema{}
}

type FederatedIdentityCredentialResourceSchema struct {
	Audience []string `tfschema:"audience"`
	Issuer   string   `tfschema:"issuer"`
	Name     string   `tfschema:"name"`

	// TODO: Remove this in V5.0
	ResourceGroupName string `tfschema:"resource_group_name,removedInNextMajorVersion"`

	ParentId               string `tfschema:"parent_id,removedInNextMajorVersion"`
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
	Subject                string `tfschema:"subject"`
}

func (r FederatedIdentityCredentialResource) Identity() resourceids.ResourceId {
	return &federatedidentitycredentials.FederatedIdentityCredentialId{}
}

func (r FederatedIdentityCredentialResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return federatedidentitycredentials.ValidateFederatedIdentityCredentialID
}

func (r FederatedIdentityCredentialResource) ResourceType() string {
	return "azurerm_federated_identity_credential"
}

func (r FederatedIdentityCredentialResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Required: true,
			ForceNew: true,
			Type:     pluginsdk.TypeString,
		},

		"user_assigned_identity_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"audience": {
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Required: true,
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
		},

		"issuer": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"subject": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
	}

	if !features.FivePointOh() {
		schema["resource_group_name"] = commonschema.ResourceGroupNameDeprecatedComputed()

		schema["parent_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			Computed:     true,
			Deprecated:   "`parent_id` has been renamed to `user_assigned_identity_id` and will be removed in v5.0 of the AzureRM Provider",
			ExactlyOneOf: []string{"user_assigned_identity_id", "parent_id"},
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		}
		schema["user_assigned_identity_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			Computed:     true,
			ExactlyOneOf: []string{"user_assigned_identity_id", "parent_id"},
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		}
	}

	return schema
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

			userAssignedIdentityId := config.UserAssignedIdentityId
			if !features.FivePointOh() && userAssignedIdentityId == "" {
				userAssignedIdentityId = config.ParentId
			}

			parentId, err := commonids.ParseUserAssignedIdentityID(userAssignedIdentityId)
			if err != nil {
				return fmt.Errorf("parsing parent resource ID: %+v", err)
			}

			locks.ByID(parentId.ID())
			defer locks.UnlockByID(parentId.ID())

			id := federatedidentitycredentials.NewFederatedIdentityCredentialID(subscriptionId, parentId.ResourceGroupName, parentId.UserAssignedIdentityName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}

func (r FederatedIdentityCredentialResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.V20241130.FederatedIdentityCredentials

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

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r FederatedIdentityCredentialResource) flatten(metadata sdk.ResourceMetaData, id *federatedidentitycredentials.FederatedIdentityCredentialId, model *federatedidentitycredentials.FederatedIdentityCredential) error {
	schema := FederatedIdentityCredentialResourceSchema{
		Name: id.FederatedIdentityCredentialName,
	}

	parentId := commonids.NewUserAssignedIdentityID(id.SubscriptionId, id.ResourceGroupName, id.UserAssignedIdentityName)
	schema.UserAssignedIdentityId = parentId.ID()

	if model != nil {
		r.mapFederatedIdentityCredentialToFederatedIdentityCredentialResourceSchema(*model, &schema)

		if !features.FivePointOh() {
			schema.ParentId = parentId.ID()
			schema.ResourceGroupName = id.ResourceGroupName
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&schema)
}

func (r FederatedIdentityCredentialResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, rmd sdk.ResourceMetaData) error {
			client := rmd.Client.ManagedIdentity.V20241130.FederatedIdentityCredentials

			var config FederatedIdentityCredentialResourceSchema
			if err := rmd.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := federatedidentitycredentials.ParseFederatedIdentityCredentialID(rmd.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil`", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			props := existing.Model.Properties

			if rmd.ResourceData.HasChange("audience") {
				props.Audiences = config.Audience
			}

			if rmd.ResourceData.HasChange("issuer") {
				props.Issuer = config.Issuer
			}

			if rmd.ResourceData.HasChange("subject") {
				props.Subject = config.Subject
			}

			parentId := commonids.NewUserAssignedIdentityID(id.SubscriptionId, id.ResourceGroupName, id.UserAssignedIdentityName)
			locks.ByID(parentId.ID())
			defer locks.UnlockByID(parentId.ID())

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
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

			userAssignedIdentityId := config.UserAssignedIdentityId
			if !features.FivePointOh() && userAssignedIdentityId == "" {
				userAssignedIdentityId = config.ParentId
			}

			parentId, err := commonids.ParseUserAssignedIdentityID(userAssignedIdentityId)
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
