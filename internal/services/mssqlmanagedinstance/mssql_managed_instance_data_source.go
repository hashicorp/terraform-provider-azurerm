// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedInstanceDataSourceModel struct {
	AdministratorLogin        string                              `tfschema:"administrator_login"`
	Collation                 string                              `tfschema:"collation"`
	CustomerManagedKeyId      string                              `tfschema:"customer_managed_key_id"`
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
	StorageSizeInGb           int                                 `tfschema:"storage_size_in_gb"`
	SubnetId                  string                              `tfschema:"subnet_id"`
	Tags                      map[string]string                   `tfschema:"tags"`
	TimezoneId                string                              `tfschema:"timezone_id"`
	VCores                    int                                 `tfschema:"vcores"`
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

			id := parse.NewManagedInstanceID(subscriptionId, state.ResourceGroupName, state.Name)
			resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := MsSqlManagedInstanceDataSourceModel{
				Name:              id.Name,
				Location:          location.NormalizeNilable(resp.Location),
				ResourceGroupName: id.ResourceGroup,
				Identity:          d.flattenIdentity(resp.Identity),
				Tags:              tags.ToTypedObject(resp.Tags),
			}

			if sku := resp.Sku; sku != nil && sku.Name != nil {
				model.SkuName = *sku.Name
			}

			if props := resp.ManagedInstanceProperties; props != nil {
				model.LicenseType = string(props.LicenseType)
				model.ProxyOverride = string(props.ProxyOverride)
				model.StorageAccountType = string(props.StorageAccountType)

				if props.AdministratorLogin != nil {
					model.AdministratorLogin = *props.AdministratorLogin
				}
				if props.Collation != nil {
					model.Collation = *props.Collation
				}
				if props.KeyID != nil {
					model.CustomerManagedKeyId = *props.KeyID
				}
				if props.FullyQualifiedDomainName != nil {
					model.Fqdn = *props.FullyQualifiedDomainName
				}
				if props.MinimalTLSVersion != nil {
					model.MinimumTlsVersion = *props.MinimalTLSVersion
				}
				if props.PublicDataEndpointEnabled != nil {
					model.PublicDataEndpointEnabled = *props.PublicDataEndpointEnabled
				}
				if props.StorageSizeInGB != nil {
					model.StorageSizeInGb = int(*props.StorageSizeInGB)
				}
				if props.SubnetID != nil {
					model.SubnetId = *props.SubnetID
				}
				if props.TimezoneID != nil {
					model.TimezoneId = *props.TimezoneID
				}
				if props.VCores != nil {
					model.VCores = int(*props.VCores)
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&model)
		},
	}
}

func (d MsSqlManagedInstanceDataSource) flattenIdentity(input *sql.ResourceIdentity) []identity.SystemOrUserAssignedList {
	if input == nil {
		return nil
	}

	principalId := ""
	if input.PrincipalID != nil {
		principalId = input.PrincipalID.String()
	}

	tenantId := ""
	if input.TenantID != nil {
		tenantId = input.TenantID.String()
	}

	var identityIds = make([]string, 0)
	for k := range input.UserAssignedIdentities {
		parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(k)
		if err != nil {
			continue
		}
		identityIds = append(identityIds, parsedId.ID())
	}

	return []identity.SystemOrUserAssignedList{{
		Type:        identity.Type(input.Type),
		PrincipalId: principalId,
		TenantId:    tenantId,
		IdentityIds: identityIds,
	}}
}
