// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
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
	GeneralPurposeV2Enabled   bool                                `tfschema:"general_purpose_v2_enabled"`
	LicenseType               string                              `tfschema:"license_type"`
	Location                  string                              `tfschema:"location"`
	MinimumTlsVersion         string                              `tfschema:"minimum_tls_version"`
	Name                      string                              `tfschema:"name"`
	PricingModel              string                              `tfschema:"pricing_model"`
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

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

func (d DataSource) ResourceType() string {
	return "azurerm_mssql_managed_instance"
}

func (d DataSource) Description() string {
	return "Use this data source to access information about an existing Microsoft SQL Azure Managed Instance."
}

func (d DataSource) ModelObject() any {
	return &MsSqlManagedInstanceDataSourceModel{}
}

func (d DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the SQL Managed Instance.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validate.ValidateMsSqlServerName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"administrator_login": {
			Description: "The administrator login name for the SQL Managed Instance.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"collation": {
			Description: "How the SQL Managed Instance will be collated.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"customer_managed_key_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"dns_zone": {
			Description: "The DNS Zone where the SQL Managed Instance is located.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"dns_zone_partner_id": {
			Description: "The ID of the SQL Managed Instance which shares the DNS zone.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"fqdn": {
			Description: "The fully qualified domain name of the Azure Managed SQL Instance.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"general_purpose_v2_enabled": {
			Description: "Whether the SQL Managed Instance is using the Next-gen General Purpose service tier.",
			Type:        schema.TypeBool,
			Computed:    true,
		},

		"license_type": {
			Description: "What type of license the SQL Managed Instance uses.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"location": commonschema.LocationComputed(),

		"minimum_tls_version": {
			Description: "The Minimum TLS Version.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"pricing_model": {
			Description: "The pricing model for the SQL Managed Instance.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"proxy_override": {
			Description: "How the SQL Managed Instance will be accessed.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"public_data_endpoint_enabled": {
			Description: "Whether the public data endpoint is enabled.",
			Type:        schema.TypeBool,
			Computed:    true,
		},

		"sku_name": {
			Description: "The SKU Name of the SQL Managed Instance.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"storage_account_type": {
			Description: "The Storage Account type used to store backups for this database.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"storage_size_in_gb": {
			Description: "Maximum storage space allocated for the SQL Managed Instance.",
			Type:        schema.TypeInt,
			Computed:    true,
		},

		"subnet_id": {
			Description: "The subnet resource ID that the SQL Managed Instance is associated with.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),

		"timezone_id": {
			Description: "The TimeZone ID that the SQL Managed Instance is running in.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		"vcores": {
			Description: "Number of cores that are assigned to the SQL Managed Instance.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
	}
}

func (d DataSource) Read() sdk.ResourceFunc {
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
				model.LicenseType = pointer.FromEnum(props.LicenseType)
				model.PricingModel = pointer.FromEnum(props.PricingModel)
				model.ProxyOverride = pointer.FromEnum(props.ProxyOverride)
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
				model.GeneralPurposeV2Enabled = pointer.From(props.IsGeneralPurposeV2)
			}

			metadata.SetID(id)
			return metadata.Encode(&model)
		},
	}
}

func (d DataSource) flattenIdentity(input *identity.LegacySystemAndUserAssignedMap) []identity.SystemOrUserAssignedList {
	if input == nil {
		return []identity.SystemOrUserAssignedList{}
	}

	identityIds := make([]string, 0)
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
