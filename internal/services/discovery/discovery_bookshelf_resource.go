// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package discovery

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
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
	KeyVaultKeyId          string `tfschema:"key_vault_key_id"`
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
}

func (r DiscoveryBookshelfResource) ResourceType() string {
	return "azurerm_discovery_bookshelf"
}

func (r DiscoveryBookshelfResource) ModelObject() any {
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
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
					},

					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
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

			if err := validateDiscoveryBookshelfCustomerManagedKeyIdentity(diff); err != nil {
				return err
			}

			if !diff.HasChange("customer_managed_key") {
				return nil
			}

			oldCMK, newCMK := diff.GetChange("customer_managed_key")
			if len(oldCMK.([]any)) != len(newCMK.([]any)) {
				return diff.ForceNew("customer_managed_key")
			}

			if !diff.NewValueKnown("customer_managed_key.0.key_vault_key_id") {
				return nil
			}

			oldKeyId, newKeyId := diff.GetChange("customer_managed_key.0.key_vault_key_id")
			oldKey, err := keyvault.ParseNestedItemID(oldKeyId.(string), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}
			newKey, err := keyvault.ParseNestedItemID(newKeyId.(string), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}

			if !strings.EqualFold(oldKey.KeyVaultBaseURL, newKey.KeyVaultBaseURL) {
				return diff.ForceNew("customer_managed_key.0.key_vault_key_id")
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
			var keyVaultProperties *bookshelves.BookshelfKeyVaultProperties
			if len(model.CustomerManagedKey) > 0 {
				customerManagedKeys = bookshelves.CustomerManagedKeysEnabled

				identityId, err := commonids.ParseUserAssignedIdentityID(model.CustomerManagedKey[0].UserAssignedIdentityId)
				if err != nil {
					return err
				}

				identity, err := metadata.Client.ManagedIdentity.V20241130.Identities.UserAssignedIdentitiesGet(ctx, *identityId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", *identityId, err)
				}

				if identity.Model == nil || identity.Model.Properties == nil || identity.Model.Properties.ClientId == nil {
					return fmt.Errorf("retrieving %s: `clientId` was nil", *identityId)
				}

				keyVaultProperties, err = expandDiscoveryBookshelfCustomerManagedKey(model.CustomerManagedKey, *identity.Model.Properties.ClientId)
				if err != nil {
					return err
				}
			}

			publicNetworkAccess := bookshelves.PublicNetworkAccessDisabled
			if model.PublicNetworkAccessEnabled {
				publicNetworkAccess = bookshelves.PublicNetworkAccessEnabled
			}

			payload := bookshelves.Bookshelf{
				Location: location.Normalize(model.Location),
				Properties: &bookshelves.BookshelfProperties{
					CustomerManagedKeys:     pointer.To(customerManagedKeys),
					KeyVaultProperties:      keyVaultProperties,
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

			state, err := flattenDiscoveryBookshelf(*id, resp.Model)
			if err != nil {
				return err
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(state)
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

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := existing.Model
			props := payload.Properties
			props.ManagedOnBehalfOfConfiguration = nil

			if metadata.ResourceData.HasChange("customer_managed_key") && props.KeyVaultProperties != nil {
				keyVaultProperties, err := expandDiscoveryBookshelfCustomerManagedKey(model.CustomerManagedKey, props.KeyVaultProperties.IdentityClientId)
				if err != nil {
					return err
				}
				props.KeyVaultProperties = keyVaultProperties
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				publicNetworkAccess := bookshelves.PublicNetworkAccessDisabled
				if model.PublicNetworkAccessEnabled {
					publicNetworkAccess = bookshelves.PublicNetworkAccessEnabled
				}
				props.PublicNetworkAccess = pointer.To(publicNetworkAccess)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
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

func flattenDiscoveryBookshelf(id bookshelves.BookshelfId, model *bookshelves.Bookshelf) (*DiscoveryBookshelfModel, error) {
	state := DiscoveryBookshelfModel{
		Name:              id.BookshelfName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		if props := model.Properties; props != nil {
			state.LogAnalyticsClusterId = pointer.From(props.LogAnalyticsClusterId)
			state.ManagedResourceGroupName = pointer.From(props.ManagedResourceGroup)
			state.PrivateEndpointSubnetId = pointer.From(props.PrivateEndpointSubnetId)
			state.PublicNetworkAccessEnabled = pointer.From(props.PublicNetworkAccess) == bookshelves.PublicNetworkAccessEnabled
			state.SearchSubnetId = pointer.From(props.SearchSubnetId)
			state.Uri = pointer.From(props.BookshelfUri)
			state.WorkloadIdentityIds = flattenDiscoveryBookshelfWorkloadIdentities(props.WorkloadIdentities)
			state.ManagedOnBehalfOfBrokerResourceIds = flattenDiscoveryBookshelfMoboBrokerResourceIds(props.ManagedOnBehalfOfConfiguration)

			if pointer.From(props.CustomerManagedKeys) == bookshelves.CustomerManagedKeysEnabled {
				customerManagedKey, err := flattenDiscoveryBookshelfCustomerManagedKey(props.KeyVaultProperties, props.WorkloadIdentities)
				if err != nil {
					return nil, err
				}
				state.CustomerManagedKey = customerManagedKey
			}
		}
	}

	return &state, nil
}

func validateDiscoveryBookshelfCustomerManagedKeyIdentity(diff *pluginsdk.ResourceDiff) error {
	identityId, ok := diff.GetOk("customer_managed_key.0.user_assigned_identity_id")
	if !ok || !diff.NewValueKnown("customer_managed_key.0.user_assigned_identity_id") {
		return nil
	}

	rawConfig := diff.GetRawConfig()
	if rawConfig.IsNull() || !rawConfig.GetAttr("workload_identity_ids").IsWhollyKnown() {
		return nil
	}

	for _, v := range diff.Get("workload_identity_ids").(*pluginsdk.Set).List() {
		if strings.EqualFold(v.(string), identityId.(string)) {
			return nil
		}
	}

	return fmt.Errorf("`customer_managed_key.0.user_assigned_identity_id` must be one of the identities specified in `workload_identity_ids`")
}

func expandDiscoveryBookshelfCustomerManagedKey(input []DiscoveryBookshelfCustomerManagedKey, identityClientId string) (*bookshelves.BookshelfKeyVaultProperties, error) {
	if len(input) == 0 {
		return nil, nil
	}

	keyId, err := keyvault.ParseNestedItemID(input[0].KeyVaultKeyId, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
	if err != nil {
		return nil, err
	}

	return &bookshelves.BookshelfKeyVaultProperties{
		IdentityClientId: identityClientId,
		KeyName:          keyId.Name,
		KeyVaultUri:      keyId.KeyVaultBaseURL,
		KeyVersion:       pointer.To(keyId.Version),
	}, nil
}

func flattenDiscoveryBookshelfCustomerManagedKey(input *bookshelves.BookshelfKeyVaultProperties, workloadIdentities *map[string]bookshelves.UserAssignedIdentity) ([]DiscoveryBookshelfCustomerManagedKey, error) {
	if input == nil {
		return []DiscoveryBookshelfCustomerManagedKey{}, nil
	}

	keyId, err := keyvault.NewNestedItemID(input.KeyVaultUri, keyvault.NestedItemTypeKey, input.KeyName, pointer.From(input.KeyVersion))
	if err != nil {
		return nil, fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
	}

	userAssignedIdentityId := ""
	for id, identity := range pointer.From(workloadIdentities) {
		if strings.EqualFold(pointer.From(identity.ClientId), input.IdentityClientId) {
			parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(id)
			if err != nil {
				return nil, err
			}
			userAssignedIdentityId = parsed.ID()
			break
		}
	}

	return []DiscoveryBookshelfCustomerManagedKey{{
		KeyVaultKeyId:          keyId.ID(),
		UserAssignedIdentityId: userAssignedIdentityId,
	}}, nil
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
