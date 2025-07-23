// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KeyVaultMHSMKeyDataSourceModel struct {
	ManagedHSMID   string                 `tfschema:"managed_hsm_id"`
	Name           string                 `tfschema:"name"`
	KeyType        string                 `tfschema:"key_type"`
	KeyOpts        []string               `tfschema:"key_opts"`
	KeySize        int64                  `tfschema:"key_size"`
	Curve          string                 `tfschema:"curve"`
	NotBeforeDate  string                 `tfschema:"not_before_date"`
	ExpirationDate string                 `tfschema:"expiration_date"`
	Tags           map[string]interface{} `tfschema:"tags"`
	VersionedId    string                 `tfschema:"versioned_id"`
	Version        string                 `tfschema:"version"`
}

type KeyvaultMHSMKeyDataSource struct{}

// Arguments implements sdk.DataSource.
func (k KeyvaultMHSMKeyDataSource) Arguments() map[string]*pluginsdk.Schema {
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
	}
}

// Attributes implements sdk.DataSource.
func (k KeyvaultMHSMKeyDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"key_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"key_size": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"curve": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"key_opts": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"not_before_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"expiration_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"versioned_id": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (k KeyvaultMHSMKeyDataSource) ModelObject() interface{} {
	return &KeyVaultMHSMKeyDataSourceModel{}
}

func (k KeyvaultMHSMKeyDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 5,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedHSMs.DataPlaneRoleDefinitionsClient
			domainSuffix, ok := metadata.Client.Account.Environment.ManagedHSM.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var config KeyVaultMHSMKeyDataSourceModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			managedHsmId, err := managedhsms.ParseManagedHSMID(config.ManagedHSMID)
			if err != nil {
				return err
			}
			id := parse.NewManagedHSMDataPlaneVersionlessKeyID(managedHsmId.ManagedHSMName, *domainSuffix, config.Name)

			resp, err := client.GetKey(ctx, id.BaseUri(), id.KeyName, "")
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					return fmt.Errorf("key %q not found", config.Name)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if key := resp.Key; key != nil {
				config.Name = id.KeyName
				config.ManagedHSMID = managedHsmId.ID()
				config.KeyType = string(key.Kty)
				config.KeyOpts = flattenKeyVaultKeyOptions(key.KeyOps)
				config.Curve = string(key.Crv)
				config.Tags = tags.Flatten(resp.Tags)

				versionedID, err := parse.ManagedHSMDataPlaneVersionedKeyID(*key.Kid, domainSuffix)
				if err != nil {
					return fmt.Errorf("parsing versioned ID: %+v", err)
				}
				config.VersionedId = versionedID.ID()
				config.Version = versionedID.KeyVersion

				if key.N != nil {
					nBytes, err := base64.RawURLEncoding.DecodeString(*key.N)
					if err != nil {
						return fmt.Errorf("could not decode N: %+v", err)
					}
					config.KeySize = int64(len(nBytes) * 8)
				}

				if attributes := resp.Attributes; attributes != nil {
					if v := attributes.NotBefore; v != nil {
						config.NotBeforeDate = time.Time(*v).Format(time.RFC3339)
					}

					if v := attributes.Expires; v != nil {
						config.ExpirationDate = time.Time(*v).Format(time.RFC3339)
					}
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&config)
		},
	}
}

func (k KeyvaultMHSMKeyDataSource) ResourceType() string {
	return "azurerm_key_vault_managed_hardware_security_module_key"
}

var _ sdk.DataSource = KeyvaultMHSMKeyDataSource{}
