// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedinstances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlManagedInstanceDataSourceModel struct {
	AdministratorLogin        string                              `tfschema:"administrator_login"`
	Collation                 string                              `tfschema:"collation"`
	CustomerManagedKeyId      string                              `tfschema:"customer_managed_key_id"`
	DnsZone                   string                              `tfschema:"dns_zone"`
	DnsZonePartnerId          string                              `tfschema:"dns_zone_partner_id"`
	Fqdn                      string                              `tfschema:"fqdn"`
	Identity                  []identity.SystemOrUserAssignedList `tfschema:"identity"`
	LicenseType               string                              `tfschema:"license_type"`
	Location                  string                              `tfschema:"location"`
	MinimumTlsVersion         string                              `tfschema:"minimum_tls_version"`
	Name                      string                              `tfschema:"name"`
	ProxyOverride             string                              `tfschema:"proxy_override"`
	PublicDataEndpointEnabled bool                                `tfschema:"public_data_endpoint_enabled"`
	ResourceGroupName         string                              `tfschema:"resource_group_name"`
	SkuName                   string                              `tfschema:"sku_name"`
	StorageAccountType        string                              `tfschema:"storage_account_type"`
	StorageSizeInGb           int64                               `tfschema:"storage_size_in_gb"`
	SubnetId                  string                              `tfschema:"subnet_id"`
	Tags                      map[string]string                   `tfschema:"tags"`
	TimezoneId                string                              `tfschema:"timezone_id"`
	VCores                    int64                               `tfschema:"vcores"`
}

var _ sdk.DataSource = MsSqlManagedInstanceDataSource{}

type MsSqlManagedInstanceDataSource struct{}

func (d MsSqlManagedInstanceDataSource) ResourceType() string {
	return "azurerm_mssql_managed_instance"
}

func (d MsSqlManagedInstanceDataSource) ModelObject() interface{} {
	return &MsSqlManagedInstanceDataSourceModel{}
}

func (d MsSqlManagedInstanceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validate.ValidateMsSqlServerName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d MsSqlManagedInstanceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"administrator_login": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"collation": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"customer_managed_key_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"dns_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"dns_zone_partner_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"fqdn": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"license_type": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"minimum_tls_version": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"proxy_override": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"public_data_endpoint_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},

		"sku_name": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"storage_size_in_gb": {
			Type:     schema.TypeInt,
			Computed: true,
		},

		"subnet_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),

		"timezone_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"vcores": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func (d MsSqlManagedInstanceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstancesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state MsSqlManagedInstanceDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := commonids.NewSqlManagedInstanceID(subscriptionId, state.ResourceGroupName, state.Name)
			resp, err := client.Get(ctx, id, managedinstances.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s model was nil", id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s properties was nil", id)
			}

			model := MsSqlManagedInstanceDataSourceModel{
				Name:              id.ManagedInstanceName,
				Location:          resp.Model.Location,
				ResourceGroupName: id.ResourceGroupName,
				Identity:          d.flattenIdentity(resp.Model.Identity),
				Tags:              pointer.From(resp.Model.Tags),
			}

			if sku := resp.Model.Sku; sku != nil {
				model.SkuName = sku.Name
			}

			if props := resp.Model.Properties; props != nil {
				model.LicenseType = string(pointer.From(props.LicenseType))
				model.ProxyOverride = string(pointer.From(props.ProxyOverride))
				model.StorageAccountType = backupStorageRedundancyToStorageAccType(pointer.From(props.RequestedBackupStorageRedundancy))
				model.AdministratorLogin = pointer.From(props.AdministratorLogin)
				model.Collation = pointer.From(props.Collation)
				model.DnsZone = pointer.From(props.DnsZone)
				model.CustomerManagedKeyId = pointer.From(props.KeyId)
				model.Fqdn = pointer.From(props.FullyQualifiedDomainName)
				model.MinimumTlsVersion = pointer.From(props.MinimalTlsVersion)
				model.PublicDataEndpointEnabled = pointer.From(props.PublicDataEndpointEnabled)
				model.StorageSizeInGb = pointer.From(props.StorageSizeInGB)
				model.SubnetId = pointer.From(props.SubnetId)
				model.TimezoneId = pointer.From(props.TimezoneId)
				model.VCores = pointer.From(props.VCores)

			}

			metadata.SetID(id)
			return metadata.Encode(&model)
		},
	}
}

func (d MsSqlManagedInstanceDataSource) flattenIdentity(input *identity.LegacySystemAndUserAssignedMap) []identity.SystemOrUserAssignedList {
	if input == nil {
		return nil
	}

	var identityIds = make([]string, 0)
	for k := range input.IdentityIds {
		parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(k)
		if err != nil {
			continue
		}
		identityIds = append(identityIds, parsedId.ID())
	}

	return []identity.SystemOrUserAssignedList{{
		Type:        input.Type,
		PrincipalId: input.PrincipalId,
		TenantId:    input.TenantId,
		IdentityIds: identityIds,
	}}
}
