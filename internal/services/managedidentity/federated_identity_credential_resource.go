package managedidentity

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2022-01-31-preview/managedidentities"
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
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
		},
		"issuer": {
			ForceNew: true,
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
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"subject": {
			ForceNew: true,
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
			client := metadata.Client.ManagedIdentity.ManagedIdentities

			var config FederatedIdentityCredentialResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			parentId, err := commonids.ParseUserAssignedIdentityID(config.ResourceName)
			if err != nil {
				return fmt.Errorf("parsing parent resource ID: %+v", err)
			}
			id := managedidentities.NewFederatedIdentityCredentialID(subscriptionId, config.ResourceGroupName, parentId.ResourceName, config.Name)

			existing, err := client.FederatedIdentityCredentialsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload managedidentities.FederatedIdentityCredential
			if err := r.mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredential(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

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
			client := metadata.Client.ManagedIdentity.ManagedIdentities
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
				schema.Name = id.FederatedIdentityCredentialResourceName
				schema.ResourceGroupName = id.ResourceGroupName
				parentId := commonids.NewUserAssignedIdentityID(id.SubscriptionId, id.ResourceGroupName, id.ResourceName)
				schema.ResourceName = parentId.ID()
				if err := r.mapFederatedIdentityCredentialToFederatedIdentityCredentialResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
func (r FederatedIdentityCredentialResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedIdentity.ManagedIdentities

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

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredentialProperties(input FederatedIdentityCredentialResourceSchema, output *managedidentities.FederatedIdentityCredentialProperties) error {

	audiences := make([]string, 0)
	for _, v := range input.Audience {
		audiences = append(audiences, v)
	}
	output.Audiences = audiences

	output.Issuer = input.Issuer
	output.Subject = input.Subject
	return nil
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialPropertiesToFederatedIdentityCredentialResourceSchema(input managedidentities.FederatedIdentityCredentialProperties, output *FederatedIdentityCredentialResourceSchema) error {

	audiences := make([]string, 0)
	for _, v := range input.Audiences {
		audiences = append(audiences, v)
	}
	output.Audience = audiences

	output.Issuer = input.Issuer
	output.Subject = input.Subject
	return nil
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredential(input FederatedIdentityCredentialResourceSchema, output *managedidentities.FederatedIdentityCredential) error {

	if output.Properties == nil {
		output.Properties = &managedidentities.FederatedIdentityCredentialProperties{}
	}
	if err := r.mapFederatedIdentityCredentialResourceSchemaToFederatedIdentityCredentialProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "FederatedIdentityCredentialProperties", "Properties", err)
	}

	return nil
}

func (r FederatedIdentityCredentialResource) mapFederatedIdentityCredentialToFederatedIdentityCredentialResourceSchema(input managedidentities.FederatedIdentityCredential, output *FederatedIdentityCredentialResourceSchema) error {

	if input.Properties == nil {
		input.Properties = &managedidentities.FederatedIdentityCredentialProperties{}
	}
	if err := r.mapFederatedIdentityCredentialPropertiesToFederatedIdentityCredentialResourceSchema(*input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "FederatedIdentityCredentialProperties", "Properties", err)
	}

	return nil
}
