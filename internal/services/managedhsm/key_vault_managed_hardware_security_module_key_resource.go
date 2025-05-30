// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

type KeyVaultMHSMKeyResource struct{}

var _ sdk.ResourceWithUpdate = KeyVaultMHSMKeyResource{}

func (r KeyVaultMHSMKeyResource) ModelObject() interface{} {
	return &KeyVaultMHSMKeyResourceSchema{}
}

type KeyVaultMHSMKeyResourceSchema struct {
	Name           string                 `tfschema:"name"`
	ManagedHSMID   string                 `tfschema:"managed_hsm_id"`
	KeyType        string                 `tfschema:"key_type"`
	KeyOpts        []string               `tfschema:"key_opts"`
	KeySize        int64                  `tfschema:"key_size"`
	Curve          string                 `tfschema:"curve"`
	NotBeforeDate  string                 `tfschema:"not_before_date"`
	ExpirationDate string                 `tfschema:"expiration_date"`
	Tags           map[string]interface{} `tfschema:"tags"`
	VersionedId    string                 `tfschema:"versioned_id"`
}

func (r KeyVaultMHSMKeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedHSMDataPlaneVersionlessKeyID
}

func (r KeyVaultMHSMKeyResource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_key"
}

func (r KeyVaultMHSMKeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"managed_hsm_id": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: managedhsms.ValidateManagedHSMID,
		},

		"key_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			// turns out Azure's *really* sensitive about the casing of these
			// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
			ValidateFunc: validation.StringInSlice([]string{
				string(keyvault.JSONWebKeyTypeECHSM),
				string(keyvault.JSONWebKeyTypeOctHSM),
				string(keyvault.JSONWebKeyTypeRSAHSM),
			}, false),
		},

		"key_size": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"curve"},
		},

		"curve": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				return old == "SECP256K1" && new == string(keyvault.JSONWebKeyCurveNameP256K)
			},
			ValidateFunc: func() pluginsdk.SchemaValidateFunc {
				out := []string{
					string(keyvault.JSONWebKeyCurveNameP256),
					string(keyvault.JSONWebKeyCurveNameP256K),
					string(keyvault.JSONWebKeyCurveNameP384),
					string(keyvault.JSONWebKeyCurveNameP521),
				}
				return validation.StringInSlice(out, false)
			}(),
			ExactlyOneOf: []string{"key_size"},
		},

		"key_opts": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				// turns out Azure's *really* sensitive about the casing of these
				// issue: https://github.com/Azure/azure-rest-api-specs/issues/1739
				ValidateFunc: validation.StringInSlice([]string{
					string(keyvault.JSONWebKeyOperationDecrypt),
					string(keyvault.JSONWebKeyOperationEncrypt),
					string(keyvault.JSONWebKeyOperationSign),
					string(keyvault.JSONWebKeyOperationUnwrapKey),
					string(keyvault.JSONWebKeyOperationVerify),
					string(keyvault.JSONWebKeyOperationWrapKey),
					string(keyvault.JSONWebKeyOperationImport),
				}, false),
			},
		},

		"not_before_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"expiration_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"tags": tags.Schema(),
	}
}

func (r KeyVaultMHSMKeyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"versioned_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (r KeyVaultMHSMKeyResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			diff := metadata.ResourceDiff

			// if any value has changed, we need to SetNewComputed on versioned_id as any change to the key is a new version
			if diff.HasChanges("key_opts", "not_before_date", "tags", "expiration_date") {
				return diff.SetNewComputed("versioned_id")
			}

			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (r KeyVaultMHSMKeyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneKeysClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config KeyVaultMHSMKeyResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managedHsmId, err := managedhsms.ParseManagedHSMID(config.ManagedHSMID)
			if err != nil {
				return err
			}
			baseUri, err := metadata.Client.ManagedHSMs.BaseUriForManagedHSM(ctx, *managedHsmId)
			if err != nil {
				return fmt.Errorf("determining the Data Plane Endpoint for %s: %+v", *managedHsmId, err)
			}
			if baseUri == nil {
				return fmt.Errorf("unable to determine the Data Plane Endpoint for %q", *managedHsmId)
			}
			endpoint, err := parse.ManagedHSMEndpoint(*baseUri, domainSuffix)
			if err != nil {
				return fmt.Errorf("parsing the Data Plane Endpoint %q: %+v", *endpoint, err)
			}

			id := parse.NewManagedHSMDataPlaneVersionlessKeyID(endpoint.ManagedHSMName, endpoint.DomainSuffix, config.Name)

			locks.ByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")
			defer locks.UnlockByName(managedHsmId.ID(), "azurerm_key_vault_managed_hardware_security_module")

			existing, err := client.GetKey(ctx, endpoint.BaseURI(), id.KeyName, "")
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := keyvault.KeyCreateParameters{
				Kty:    keyvault.JSONWebKeyType(config.KeyType),
				KeyOps: expandKeyVaultKeyOptions(config.KeyOpts),
				KeyAttributes: &keyvault.KeyAttributes{
					Enabled: utils.Bool(true),
				},

				Tags: tags.Expand(config.Tags),
			}

			if config.Curve != "" {
				if config.KeyType != string(keyvault.JSONWebKeyTypeECHSM) {
					return fmt.Errorf("`key_type` must be `EC-HSM` when `curve` is set")
				}
				parameters.Curve = keyvault.JSONWebKeyCurveName(config.Curve)
			}

			if config.KeySize > 0 {
				if config.KeyType != string(keyvault.JSONWebKeyTypeRSAHSM) {
					return fmt.Errorf("`key_type` must be `RSA-HSM` when `key_size` is set")
				}
				parameters.KeySize = pointer.To(int32(config.KeySize))
			}

			if config.NotBeforeDate != "" {
				notBeforeDate, _ := time.Parse(time.RFC3339, config.NotBeforeDate) // validated by schema
				notBeforeUnixTime := date.UnixTime(notBeforeDate)
				parameters.KeyAttributes.NotBefore = &notBeforeUnixTime
			}

			if config.ExpirationDate != "" {
				expirationDate, _ := time.Parse(time.RFC3339, config.ExpirationDate) // validated by schema
				expirationUnixTime := date.UnixTime(expirationDate)
				parameters.KeyAttributes.Expires = &expirationUnixTime
			}

			if resp, err := client.CreateKey(ctx, endpoint.BaseURI(), config.Name, parameters); err != nil {
				if metadata.Client.Features.KeyVault.RecoverSoftDeletedHSMKeys && utils.ResponseWasConflict(resp.Response) {
					recoveredKey, err := client.RecoverDeletedKey(ctx, endpoint.BaseURI(), config.Name)
					if err != nil {
						return err
					}
					log.Printf("[DEBUG] Recovering HSM Key %q with ID: %q", config.Name, *recoveredKey.Key.Kid)
					if kid := recoveredKey.Key.Kid; kid != nil {
						stateConf := &pluginsdk.StateChangeConf{
							Pending:                   []string{"pending"},
							Target:                    []string{"available"},
							Refresh:                   managedHSMKeyRefreshFunc(*kid),
							Delay:                     30 * time.Second,
							PollInterval:              10 * time.Second,
							ContinuousTargetOccurence: 10,
							Timeout:                   metadata.ResourceData.Timeout(pluginsdk.TimeoutCreate),
						}

						if _, err := stateConf.WaitForStateContext(ctx); err != nil {
							return fmt.Errorf("waiting for HSM Key %q to become available: %s", config.Name, err)
						}
						log.Printf("[DEBUG] Key %q recovered with ID: %q", config.Name, *kid)
					}
				} else {
					return fmt.Errorf("Creating Key: %+v", err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KeyVaultMHSMKeyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneKeysClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			schema := KeyVaultMHSMKeyResourceSchema{}

			id, err := parse.ManagedHSMDataPlaneVersionlessKeyID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			resourceManagerId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseUri(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining Resource Manager ID for %q: %+v", id, err)
			}
			if resourceManagerId == nil {
				return metadata.MarkAsGone(*id)
			}

			resp, err := client.GetKey(ctx, id.BaseUri(), id.KeyName, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if key := resp.Key; key != nil {
				schema.Name = id.KeyName
				schema.ManagedHSMID = resourceManagerId.ID()
				schema.KeyType = string(key.Kty)
				schema.KeyOpts = flattenKeyVaultKeyOptions(key.KeyOps)
				schema.Curve = string(key.Crv)
				schema.Tags = tags.Flatten(resp.Tags)
				schema.VersionedId = pointer.From(key.Kid)
				if key.N != nil {
					nBytes, err := base64.RawURLEncoding.DecodeString(*key.N)
					if err != nil {
						return fmt.Errorf("Could not decode N: %+v", err)
					}
					schema.KeySize = int64(len(nBytes) * 8)
				}

				if attributes := resp.Attributes; attributes != nil {
					if v := attributes.NotBefore; v != nil {
						schema.NotBeforeDate = time.Time(*v).Format(time.RFC3339)
					}

					if v := attributes.Expires; v != nil {
						schema.ExpirationDate = time.Time(*v).Format(time.RFC3339)
					}
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r KeyVaultMHSMKeyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleAssignmentsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config KeyVaultMHSMKeyResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := parse.ManagedHSMDataPlaneVersionlessKeyID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			resourceManagerId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseUri(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining Resource Manager ID for %q: %+v", id, err)
			}
			if resourceManagerId == nil {
				return fmt.Errorf("unable to determine the Resource Manager ID for %s", id)
			}

			parameters := keyvault.KeyUpdateParameters{
				KeyOps: expandKeyVaultKeyOptions(config.KeyOpts),
				KeyAttributes: &keyvault.KeyAttributes{
					Enabled: utils.Bool(true),
				},

				Tags: tags.Expand(config.Tags),
			}

			if config.NotBeforeDate != "" {
				notBeforeDate, _ := time.Parse(time.RFC3339, config.NotBeforeDate) // validated by schema
				notBeforeUnixTime := date.UnixTime(notBeforeDate)
				parameters.KeyAttributes.NotBefore = &notBeforeUnixTime
			}

			if config.ExpirationDate != "" {
				expirationDate, _ := time.Parse(time.RFC3339, config.ExpirationDate) // validated by schema
				expirationUnixTime := date.UnixTime(expirationDate)
				parameters.KeyAttributes.Expires = &expirationUnixTime
			}

			if _, err = client.UpdateKey(ctx, id.BaseUri(), config.Name, "", parameters); err != nil {
				return err
			}

			return nil
		},
	}
}

var _ deleteAndPurgeNestedItem = deleteAndPurgeKey{}

type deleteAndPurgeKey struct {
	client      *keyvault.BaseClient
	keyVaultUri string
	name        string
}

func (d deleteAndPurgeKey) DeleteNestedItem(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.DeleteKey(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func (d deleteAndPurgeKey) NestedItemHasBeenDeleted(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetKey(ctx, d.keyVaultUri, d.name, "")
	return resp.Response, err
}

func (d deleteAndPurgeKey) PurgeNestedItem(ctx context.Context) (autorest.Response, error) {
	return d.client.PurgeDeletedKey(ctx, d.keyVaultUri, d.name)
}

func (d deleteAndPurgeKey) NestedItemHasBeenPurged(ctx context.Context) (autorest.Response, error) {
	resp, err := d.client.GetDeletedKey(ctx, d.keyVaultUri, d.name)
	return resp.Response, err
}

func (r KeyVaultMHSMKeyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneKeysClient
			hsmClient := metadata.Client.ManagedHSMs.ManagedHsmClient

			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := parse.ManagedHSMDataPlaneVersionlessKeyID(metadata.ResourceData.Id(), domainSuffix)
			if err != nil {
				return err
			}

			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)
			resourceManagerId, err := metadata.Client.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseUri(), domainSuffix)
			if err != nil {
				return fmt.Errorf("determining Resource Manager ID for %q: %+v", id, err)
			}
			if resourceManagerId == nil {
				return fmt.Errorf("unable to determine the Resource Manager ID for %s", id)
			}

			managedHSM, err := hsmClient.Get(ctx, *resourceManagerId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", resourceManagerId, err)
			}

			shouldPurge := metadata.Client.Features.KeyVault.PurgeSoftDeletedHSMKeysOnDestroy
			if shouldPurge && managedHSM.Model != nil && managedHSM.Model.Properties != nil && pointer.From(managedHSM.Model.Properties.EnablePurgeProtection) {
				log.Printf("[DEBUG] cannot purge key %q because Managed HSM %q has purge protection enabled", id.KeyName, id.ManagedHSMName)
				shouldPurge = false
			}

			description := fmt.Sprintf("Key %q (Managed HSM %q)", id.KeyName, id.ManagedHSMName)
			deleter := deleteAndPurgeKey{
				client:      client,
				keyVaultUri: id.BaseUri(),
				name:        id.KeyName,
			}

			return deleteAndOptionallyPurge(ctx, description, shouldPurge, deleter)
		},
	}
}

func expandKeyVaultKeyOptions(input []string) *[]keyvault.JSONWebKeyOperation {
	results := make([]keyvault.JSONWebKeyOperation, 0, len(input))

	for _, option := range input {
		results = append(results, keyvault.JSONWebKeyOperation(option))
	}

	return &results
}

func flattenKeyVaultKeyOptions(input *[]string) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	return append(results, *input...)
}
