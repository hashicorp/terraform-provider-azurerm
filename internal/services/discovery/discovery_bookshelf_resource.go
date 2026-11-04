// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package discovery

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity

var (
	_ sdk.ResourceWithUpdate        = DiscoveryBookshelfResource{}
	_ sdk.ResourceWithIdentity      = DiscoveryBookshelfResource{}
	_ sdk.ResourceWithCustomizeDiff = DiscoveryBookshelfResource{}
)

type DiscoveryBookshelfResource struct{}

type DiscoveryBookshelfModel struct {
	Name                       string                                 `tfschema:"name"`
	ResourceGroupName          string                                 `tfschema:"resource_group_name"`
	Location                   string                                 `tfschema:"location"`
	CustomerManagedKey         []DiscoveryBookshelfCustomerManagedKey `tfschema:"customer_managed_key"`
	LogAnalyticsClusterId      string                                 `tfschema:"log_analytics_cluster_id"`
	PrivateEndpointSubnetId    string                                 `tfschema:"private_endpoint_subnet_id"`
	PublicNetworkAccessEnabled bool                                   `tfschema:"public_network_access_enabled"`
	SearchSubnetId             string                                 `tfschema:"search_subnet_id"`
	WorkloadIdentityIds        []string                               `tfschema:"workload_identity_ids"`
	Tags                       map[string]string                      `tfschema:"tags"`

	ManagedOnBehalfOfBrokerResourceIds []string `tfschema:"managed_on_behalf_of_broker_resource_ids"`
	ManagedResourceGroupName           string   `tfschema:"managed_resource_group_name"`
	Uri                                string   `tfschema:"uri"`
}

type DiscoveryBookshelfCustomerManagedKey struct {
	IdentityClientId string `tfschema:"identity_client_id"`
	KeyName          string `tfschema:"key_name"`
	KeyVaultUri      string `tfschema:"key_vault_uri"`
	KeyVersion       string `tfschema:"key_version"`
}

func (r DiscoveryBookshelfResource) ResourceType() string {
	return "azurerm_discovery_bookshelf"
}

func (r DiscoveryBookshelfResource) ModelObject() interface{} {
	return &DiscoveryBookshelfModel{}
}

func (r DiscoveryBookshelfResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return bookshelves.ValidateBookshelfID
}

func (r DiscoveryBookshelfResource) Identity() resourceids.ResourceId {
	return &bookshelves.BookshelfId{}
}

func (r DiscoveryBookshelfResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9-]{3,24}$`),
				"`name` must be between 3 and 24 characters long and may only contain letters, numbers and hyphens",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"customer_managed_key": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			RequiredWith: []string{"log_analytics_cluster_id"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity_client_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsUUID,
					},

					"key_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"key_vault_uri": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},

					"key_version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"log_analytics_cluster_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: clusters.ValidateClusterID,
		},

		"private_endpoint_subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"search_subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"workload_identity_ids": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r DiscoveryBookshelfResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_on_behalf_of_broker_resource_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"managed_resource_group_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DiscoveryBookshelfResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff

			if !diff.HasChange("customer_managed_key") {
				return nil
			}

			oldCMK, newCMK := diff.GetChange("customer_managed_key")
			if len(oldCMK.([]interface{})) != len(newCMK.([]interface{})) {
				return diff.ForceNew("customer_managed_key")
			}

			oldVersion, newVersion := diff.GetChange("customer_managed_key.0.key_version")
			if oldVersion.(string) != "" && newVersion.(string) == "" {
				return diff.ForceNew("customer_managed_key.0.key_version")
			}

			return nil
		},
	}
}

func (r DiscoveryBookshelfResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.Bookshelves
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DiscoveryBookshelfModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := bookshelves.NewBookshelfID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			customerManagedKeys := bookshelves.CustomerManagedKeysDisabled
			if len(model.CustomerManagedKey) > 0 {
				customerManagedKeys = bookshelves.CustomerManagedKeysEnabled
			}

			publicNetworkAccess := bookshelves.PublicNetworkAccessDisabled
			if model.PublicNetworkAccessEnabled {
				publicNetworkAccess = bookshelves.PublicNetworkAccessEnabled
			}

			payload := bookshelves.Bookshelf{
				Location: location.Normalize(model.Location),
				Properties: &bookshelves.BookshelfProperties{
					CustomerManagedKeys:     pointer.To(customerManagedKeys),
					KeyVaultProperties:      expandDiscoveryBookshelfCustomerManagedKey(model.CustomerManagedKey),
					PrivateEndpointSubnetId: pointer.To(model.PrivateEndpointSubnetId),
					PublicNetworkAccess:     pointer.To(publicNetworkAccess),
					SearchSubnetId:          pointer.To(model.SearchSubnetId),
					WorkloadIdentities:      expandDiscoveryBookshelfWorkloadIdentities(model.WorkloadIdentityIds),
				},
				Tags: pointer.To(model.Tags),
			}

			if model.LogAnalyticsClusterId != "" {
				payload.Properties.LogAnalyticsClusterId = pointer.To(model.LogAnalyticsClusterId)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r DiscoveryBookshelfResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.Bookshelves

			id, err := bookshelves.ParseBookshelfID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			state := DiscoveryBookshelfModel{
				Name:              id.BookshelfName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(resp.Model.Location),
				Tags:              pointer.From(resp.Model.Tags),
			}

			if props := resp.Model.Properties; props != nil {
				state.LogAnalyticsClusterId = pointer.From(props.LogAnalyticsClusterId)
				state.ManagedResourceGroupName = pointer.From(props.ManagedResourceGroup)
				state.PrivateEndpointSubnetId = pointer.From(props.PrivateEndpointSubnetId)
				state.PublicNetworkAccessEnabled = pointer.From(props.PublicNetworkAccess) == bookshelves.PublicNetworkAccessEnabled
				state.SearchSubnetId = pointer.From(props.SearchSubnetId)
				state.Uri = pointer.From(props.BookshelfUri)
				state.WorkloadIdentityIds = flattenDiscoveryBookshelfWorkloadIdentities(props.WorkloadIdentities)
				state.ManagedOnBehalfOfBrokerResourceIds = flattenDiscoveryBookshelfMoboBrokerResourceIds(props.ManagedOnBehalfOfConfiguration)

				if pointer.From(props.CustomerManagedKeys) == bookshelves.CustomerManagedKeysEnabled {
					state.CustomerManagedKey = flattenDiscoveryBookshelfCustomerManagedKey(props.KeyVaultProperties)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DiscoveryBookshelfResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.Bookshelves

			id, err := bookshelves.ParseBookshelfID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DiscoveryBookshelfModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := bookshelves.BookshelfUpdate{
				Properties: &bookshelves.BookshelfPropertiesUpdate{},
			}

			if metadata.ResourceData.HasChange("customer_managed_key") && len(model.CustomerManagedKey) > 0 {
				payload.Properties.KeyVaultProperties = &bookshelves.BookshelfKeyVaultPropertiesUpdate{
					KeyName: pointer.To(model.CustomerManagedKey[0].KeyName),
				}

				if model.CustomerManagedKey[0].KeyVersion != "" {
					payload.Properties.KeyVaultProperties.KeyVersion = pointer.To(model.CustomerManagedKey[0].KeyVersion)
				}
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				publicNetworkAccess := bookshelves.PublicNetworkAccessDisabled
				if model.PublicNetworkAccessEnabled {
					publicNetworkAccess = bookshelves.PublicNetworkAccessEnabled
				}
				payload.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DiscoveryBookshelfResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.Bookshelves

			id, err := bookshelves.ParseBookshelfID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDiscoveryBookshelfCustomerManagedKey(input []DiscoveryBookshelfCustomerManagedKey) *bookshelves.BookshelfKeyVaultProperties {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	output := bookshelves.BookshelfKeyVaultProperties{
		IdentityClientId: v.IdentityClientId,
		KeyName:          v.KeyName,
		KeyVaultUri:      v.KeyVaultUri,
	}

	if v.KeyVersion != "" {
		output.KeyVersion = pointer.To(v.KeyVersion)
	}

	return &output
}

func flattenDiscoveryBookshelfCustomerManagedKey(input *bookshelves.BookshelfKeyVaultProperties) []DiscoveryBookshelfCustomerManagedKey {
	if input == nil {
		return []DiscoveryBookshelfCustomerManagedKey{}
	}

	return []DiscoveryBookshelfCustomerManagedKey{{
		IdentityClientId: input.IdentityClientId,
		KeyName:          input.KeyName,
		KeyVaultUri:      input.KeyVaultUri,
		KeyVersion:       pointer.From(input.KeyVersion),
	}}
}

func expandDiscoveryBookshelfWorkloadIdentities(input []string) *map[string]bookshelves.UserAssignedIdentity {
	if len(input) == 0 {
		return nil
	}

	output := make(map[string]bookshelves.UserAssignedIdentity, len(input))
	for _, id := range input {
		output[id] = bookshelves.UserAssignedIdentity{}
	}

	return &output
}

func flattenDiscoveryBookshelfWorkloadIdentities(input *map[string]bookshelves.UserAssignedIdentity) []string {
	output := make([]string, 0)
	if input == nil {
		return output
	}

	for id := range *input {
		if parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(id); err == nil {
			output = append(output, parsed.ID())
			continue
		}
		output = append(output, id)
	}

	return output
}

func flattenDiscoveryBookshelfMoboBrokerResourceIds(input *bookshelves.WithMoboBrokerResources) []string {
	output := make([]string, 0)
	if input == nil || input.MoboBrokerResources == nil {
		return output
	}

	for _, v := range *input.MoboBrokerResources {
		output = append(output, pointer.From(v.Id))
	}

	return output
}
