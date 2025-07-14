// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfiguration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

type KeyResource struct{}

var _ sdk.ResourceWithCustomizeDiff = KeyResource{}

var _ sdk.ResourceWithStateMigration = KeyResource{}

const (
	KeyTypeVault        = "vault"
	KeyTypeKV           = "kv"
	VaultKeyContentType = "application/vnd.microsoft.appconfig.keyvaultref+json;charset=utf-8"
)

type KeyResourceModel struct {
	ConfigurationStoreId string                 `tfschema:"configuration_store_id"`
	Key                  string                 `tfschema:"key"`
	ContentType          string                 `tfschema:"content_type"`
	Etag                 string                 `tfschema:"etag"`
	Label                string                 `tfschema:"label"`
	Value                string                 `tfschema:"value"`
	Locked               bool                   `tfschema:"locked"`
	Tags                 map[string]interface{} `tfschema:"tags"`
	Type                 string                 `tfschema:"type"`
	VaultKeyReference    string                 `tfschema:"vault_key_reference"`
}

type VaultKeyReference struct {
	URI string `json:"uri"`
}

func (k KeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"configuration_store_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			// User-specified segments are lowercased in the API response
			// tracked in https://github.com/Azure/azure-rest-api-specs/issues/24337
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     configurationstores.ValidateConfigurationStoreID,
		},
		"key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"content_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C We set some values in this field depending on the `type` so this needs to remain Computed
			Computed: true,
		},
		"etag": {
			Type: pluginsdk.TypeString,
			// NOTE: O+C The value of this is updated anytime the resource changes so this should remain Computed
			Computed: true,
			Optional: true,
		},
		"label": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"value": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ConflictsWith: []string{
				"vault_key_reference",
			},
			// if `type` is set to `vault`, then `value` will be set by `vault_key_reference`
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				return d.Get("type").(string) == KeyTypeVault && d.Get("vault_key_reference").(string) != "" && old != "" && new == ""
			},
		},
		"locked": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      KeyTypeKV,
			ValidateFunc: validation.StringInSlice([]string{KeyTypeVault, KeyTypeKV}, false),
		},
		"vault_key_reference": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			ConflictsWith: []string{
				"value",
			},
		},
		"tags": tags.Schema(),
	}
}

func (k KeyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (k KeyResource) ModelObject() interface{} {
	return &KeyResourceModel{}
}

func (k KeyResource) ResourceType() string {
	return "azurerm_app_configuration_key"
}

func (k KeyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model KeyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			configurationStoreId, err := configurationstores.ParseConfigurationStoreID(model.ConfigurationStoreId)
			if err != nil {
				return err
			}

			configurationStoreEndpoint, err := metadata.Client.AppConfiguration.EndpointForConfigurationStore(ctx, *configurationStoreId)
			if err != nil {
				return fmt.Errorf("retrieving Endpoint for feature %q in %q: %s", model.Key, *configurationStoreId, err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(*configurationStoreEndpoint)
			if err != nil {
				return err
			}

			nestedItemId, err := parse.NewNestedItemID(client.Endpoint, model.Key, model.Label)
			if err != nil {
				return err
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return errors.New("internal-error: context had no deadline")
			}

			// from https://learn.microsoft.com/en-us/azure/azure-app-configuration/concept-enable-rbac#azure-built-in-roles-for-azure-app-configuration
			// allow some time for role permission to be propagated
			metadata.Logger.Infof("[DEBUG] Waiting for App Configuration Key %q read permission to be propagated", model.Key)
			stateConf := &pluginsdk.StateChangeConf{
				Pending:                   []string{"Forbidden"},
				Target:                    []string{"Error", "Exists", "NotFound"},
				Refresh:                   appConfigurationGetKeyRefreshFunc(ctx, client, model.Key, model.Label),
				PollInterval:              10 * time.Second,
				ContinuousTargetOccurence: 3,
				Timeout:                   time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for App Configuration Key %q read permission to be propagated: %+v", model.Key, err)
			}

			kv, err := client.GetKeyValue(ctx, model.Key, model.Label, "", "", "", []appconfiguration.KeyValueFields{})
			if err != nil {
				if v, ok := err.(autorest.DetailedError); ok {
					if !utils.ResponseWasNotFound(autorest.Response{Response: v.Response}) {
						return fmt.Errorf("checking for presence of existing %s: %+v", nestedItemId, err)
					}
				} else {
					return fmt.Errorf("while checking for key's %q existence: %+v", model.Key, err)
				}
			} else if kv.Response.StatusCode == 200 {
				return tf.ImportAsExistsError(k.ResourceType(), nestedItemId.ID())
			}

			entity := appconfiguration.KeyValue{
				Key:   pointer.To(model.Key),
				Label: pointer.To(model.Label),
				Tags:  tags.Expand(model.Tags),
			}

			switch model.Type {
			case KeyTypeKV:
				entity.ContentType = pointer.To(model.ContentType)
				entity.Value = pointer.To(model.Value)
			case KeyTypeVault:
				entity.ContentType = pointer.To(VaultKeyContentType)
				ref, err := json.Marshal(VaultKeyReference{URI: model.VaultKeyReference})
				if err != nil {
					return fmt.Errorf("while encoding vault key reference: %+v", err)
				}
				entity.Value = pointer.To(string(ref))
			}

			if _, err = client.PutKeyValue(ctx, model.Key, model.Label, &entity, "", ""); err != nil {
				return err
			}

			if model.Locked {
				_, err = client.PutLock(ctx, model.Key, model.Label, "", "")
				if err != nil {
					return fmt.Errorf("while locking key/label pair %q/%q: %+v", model.Key, model.Label, err)
				}
			}

			// https://github.com/Azure/AppConfiguration/issues/763
			metadata.Logger.Infof("[DEBUG] Waiting for App Configuration Key %q to be provisioned", model.Key)
			stateConf = &pluginsdk.StateChangeConf{
				Pending:                   []string{"NotFound", "Forbidden"},
				Target:                    []string{"Exists"},
				Refresh:                   appConfigurationGetKeyRefreshFunc(ctx, client, model.Key, model.Label),
				PollInterval:              5 * time.Second,
				ContinuousTargetOccurence: 4,
				Timeout:                   time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for App Configuration Key %q to be provisioned: %+v", model.Key, err)
			}

			metadata.SetID(nestedItemId)
			return nil
		},
		Timeout: 45 * time.Minute,
	}
}

func (k KeyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			nestedItemId, err := parse.ParseNestedItemID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			domainSuffix, ok := metadata.Client.Account.Environment.AppConfiguration.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine AppConfiguration domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			configurationStoreIdRaw, err := metadata.Client.AppConfiguration.ConfigurationStoreIDFromEndpoint(ctx, subscriptionId, nestedItemId.ConfigurationStoreEndpoint, *domainSuffix)
			if err != nil {
				return fmt.Errorf("while retrieving the Resource ID of Configuration Store at Endpoint: %q: %s", nestedItemId.ConfigurationStoreEndpoint, err)
			}
			if configurationStoreIdRaw == nil {
				// if the AppConfiguration is gone then all the data inside it is too
				log.Printf("[DEBUG] Unable to determine the Resource ID for Configuration Store at Endpoint %q - removing from state", nestedItemId.ConfigurationStoreEndpoint)
				return metadata.MarkAsGone(nestedItemId)
			}

			configurationStoreId, err := configurationstores.ParseConfigurationStoreID(*configurationStoreIdRaw)
			if err != nil {
				return err
			}

			exists, err := metadata.Client.AppConfiguration.Exists(ctx, *configurationStoreId)
			if err != nil {
				return fmt.Errorf("while checking Configuration Store %q for feature %q existence: %v", *configurationStoreId, *nestedItemId, err)
			}
			if !exists {
				log.Printf("[DEBUG] Configuration Store %q for feature %q was not found - removing from state", *configurationStoreId, *nestedItemId)
				return metadata.MarkAsGone(nestedItemId)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(nestedItemId.ConfigurationStoreEndpoint)
			if err != nil {
				return err
			}

			kv, err := client.GetKeyValue(ctx, nestedItemId.Key, nestedItemId.Label, "", "", "", []appconfiguration.KeyValueFields{})
			if err != nil {
				if v, ok := err.(autorest.DetailedError); ok {
					if utils.ResponseWasNotFound(autorest.Response{Response: v.Response}) {
						return metadata.MarkAsGone(nestedItemId)
					}
				} else {
					return fmt.Errorf("while checking for key %q existence: %+v", *nestedItemId, err)
				}
				return fmt.Errorf("while checking for key %q existence: %+v", *nestedItemId, err)
			}

			model := KeyResourceModel{
				ConfigurationStoreId: configurationStoreId.ID(),
				Key:                  pointer.From(kv.Key),
				ContentType:          pointer.From(kv.ContentType),
				Etag:                 pointer.From(kv.Etag),
				Label:                pointer.From(kv.Label),
				Tags:                 tags.Flatten(kv.Tags),
			}

			if pointer.From(kv.ContentType) != VaultKeyContentType {
				model.Type = KeyTypeKV
				model.Value = pointer.From(kv.Value)
			} else {
				var ref VaultKeyReference
				refBytes := []byte(pointer.From(kv.Value))
				err := json.Unmarshal(refBytes, &ref)
				if err != nil {
					return fmt.Errorf("while unmarshalling vault reference: %+v", err)
				}

				model.Type = KeyTypeVault
				model.VaultKeyReference = ref.URI
				model.ContentType = VaultKeyContentType
				model.Value = pointer.From(kv.Value)
			}

			if kv.Locked != nil {
				model.Locked = *kv.Locked
			}
			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

func (k KeyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			nestedItemId, err := parse.ParseNestedItemID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(nestedItemId.ConfigurationStoreEndpoint)
			if err != nil {
				return err
			}

			var model KeyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			configurationStoreId, err := configurationstores.ParseConfigurationStoreID(model.ConfigurationStoreId)
			if err != nil {
				return err
			}

			metadata.Client.AppConfiguration.AddToCache(*configurationStoreId, nestedItemId.ConfigurationStoreEndpoint)

			if metadata.ResourceData.HasChange("value") || metadata.ResourceData.HasChange("content_type") || metadata.ResourceData.HasChange("tags") || metadata.ResourceData.HasChange("type") || metadata.ResourceData.HasChange("vault_key_reference") {
				entity := appconfiguration.KeyValue{
					Key:   pointer.To(model.Key),
					Label: pointer.To(model.Label),
					Tags:  tags.Expand(model.Tags),
				}

				switch model.Type {
				case KeyTypeKV:
					entity.ContentType = pointer.To(model.ContentType)
					entity.Value = pointer.To(model.Value)
				case KeyTypeVault:
					entity.ContentType = pointer.To(VaultKeyContentType)
					ref, err := json.Marshal(VaultKeyReference{URI: model.VaultKeyReference})
					if err != nil {
						return fmt.Errorf("while encoding vault key reference: %+v", err)
					}
					entity.Value = pointer.To(string(ref))
				}
				if _, err = client.PutKeyValue(ctx, model.Key, model.Label, &entity, "", ""); err != nil {
					return fmt.Errorf("while updating key/label pair %s/%s: %+v", model.Key, model.Label, err)
				}
			}

			if metadata.ResourceData.HasChange("locked") {
				if model.Locked {
					if _, err = client.PutLock(ctx, model.Key, model.Label, "", ""); err != nil {
						return fmt.Errorf("while locking key/label pair %s/%s: %+v", model.Key, model.Label, err)
					}
				} else {
					if _, err = client.DeleteLock(ctx, model.Key, model.Label, "", ""); err != nil {
						return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", model.Key, model.Label, err)
					}
				}
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k KeyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			nestedItemId, err := parse.ParseNestedItemID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(nestedItemId.ConfigurationStoreEndpoint)
			if err != nil {
				return err
			}

			if _, err = client.DeleteLock(ctx, nestedItemId.Key, nestedItemId.Label, "", ""); err != nil {
				return fmt.Errorf("while unlocking key %q: %+v", *nestedItemId, err)
			}

			if _, err = client.DeleteKeyValue(ctx, nestedItemId.Key, nestedItemId.Label, ""); err != nil {
				return fmt.Errorf("while removing key %q: %+v", *nestedItemId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k KeyResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			keyType := rd.Get("type").(string)

			if keyType == KeyTypeVault {
				contentType := rd.Get("content_type").(string)
				if rd.HasChange("content_type") && contentType != VaultKeyContentType {
					return fmt.Errorf("key type %q cannot have content type other than %q (found %q)", KeyTypeVault, VaultKeyContentType, contentType)
				}

				if rd.HasChange("value") && rd.Get("value").(string) != "" {
					return fmt.Errorf("'value' should only be set when key type is set to %q", KeyTypeKV)
				}
			}

			if keyType == KeyTypeKV && rd.Get("vault_key_reference").(string) != "" {
				return fmt.Errorf("'vault_key_reference' should only be set when key type is set to %q", KeyTypeVault)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (k KeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NestedItemId
}

func (k KeyResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 2,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.KeyResourceV0ToV1{},
			1: migration.KeyResourceV1ToV2{},
		},
	}
}
