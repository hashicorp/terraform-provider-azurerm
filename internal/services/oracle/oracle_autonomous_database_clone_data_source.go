// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = AutonomousDatabaseCloneDataSource{}

type AutonomousDatabaseCloneDataSource struct{}

type AutonomousDatabaseCloneDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`

	// Clone-specific properties
	SourceId                string `tfschema:"source_id"`
	CloneType               string `tfschema:"clone_type"`
	IsReconnectCloneEnabled bool   `tfschema:"is_reconnect_clone_enabled"`
	IsRefreshableClone      bool   `tfschema:"is_refreshable_clone"`
	RefreshableModel        string `tfschema:"refreshable_model"`
	RefreshableStatus       string `tfschema:"refreshable_status"`
	TimeUntilReconnectClone string `tfschema:"time_until_reconnect_clone"`

	// Base properties (computed)
	AutonomousDatabaseId         string   `tfschema:"autonomous_database_id"`
	BackupRetentionPeriodInDays  int64    `tfschema:"backup_retention_period_in_days"`
	CharacterSet                 string   `tfschema:"character_set"`
	ComputeCount                 float64  `tfschema:"compute_count"`
	ComputeModel                 string   `tfschema:"compute_model"`
	ConnectionStrings            []string `tfschema:"connection_strings"`
	CustomerContacts             []string `tfschema:"customer_contacts"`
	DataStorageSizeInGbs         int64    `tfschema:"data_storage_size_in_gbs"`
	DataStorageSizeInTbs         int64    `tfschema:"data_storage_size_in_tbs"`
	DbVersion                    string   `tfschema:"db_version"`
	DbWorkload                   string   `tfschema:"db_workload"`
	DisplayName                  string   `tfschema:"display_name"`
	LicenseModel                 string   `tfschema:"license_model"`
	AutoScalingEnabled           bool     `tfschema:"auto_scaling_enabled"`
	AutoScalingForStorageEnabled bool     `tfschema:"auto_scaling_for_storage_enabled"`
	MtlsConnectionRequired       bool     `tfschema:"mtls_connection_required"`
	NationalCharacterSet         string   `tfschema:"national_character_set"`
	SubnetId                     string   `tfschema:"subnet_id"`
	VnetId                       string   `tfschema:"virtual_network_id"`
	LifecycleState               string   `tfschema:"lifecycle_state"`
	PrivateEndpoint              string   `tfschema:"private_endpoint"`
	PrivateEndpointIp            string   `tfschema:"private_endpoint_ip"`
	ServiceConsoleUrl            string   `tfschema:"service_console_url"`
	SqlWebDeveloperUrl           string   `tfschema:"sql_web_developer_url"`
	TimeCreated                  string   `tfschema:"time_created"`
	OciUrl                       string   `tfschema:"oci_url"`
}

func (AutonomousDatabaseCloneDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (AutonomousDatabaseCloneDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),

		// Clone-specific properties
		"source_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"clone_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"is_reconnect_clone_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"is_refreshable_clone": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"refreshable_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"refreshable_status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_until_reconnect_clone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		// Base properties
		"autonomous_database_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"backup_retention_period_in_days": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"character_set": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"compute_count": {
			Type:     pluginsdk.TypeFloat,
			Computed: true,
		},

		"compute_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"connection_strings": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"customer_contacts": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"data_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"data_storage_size_in_tbs": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"db_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"db_workload": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"license_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"auto_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"auto_scaling_for_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"mtls_connection_required": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"national_character_set": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"virtual_network_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"lifecycle_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_endpoint_ip": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"service_console_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"sql_web_developer_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_created": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"oci_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (AutonomousDatabaseCloneDataSource) ModelObject() interface{} {
	return &AutonomousDatabaseCloneDataSourceModel{}
}

func (AutonomousDatabaseCloneDataSource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_clone"
}

func (AutonomousDatabaseCloneDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Oracle.OracleClient.AutonomousDatabases

			var model AutonomousDatabaseCloneDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := autonomousdatabases.NewAutonomousDatabaseID(subscriptionId,
				model.ResourceGroupName,
				model.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model.Name = id.AutonomousDatabaseName
			model.ResourceGroupName = id.ResourceGroupName

			if m := resp.Model; m != nil {
				model.Location = m.Location
				model.Tags = pointer.From(m.Tags)

				// Verify this is actually a clone database
				props, ok := m.Properties.(autonomousdatabases.AutonomousDatabaseCloneProperties)
				if !ok {
					return fmt.Errorf("%s is not a clone type autonomous database", id)
				}

				// Clone-specific properties
				model.CloneType = string(props.CloneType)
				model.SourceId = props.SourceId
				model.IsReconnectCloneEnabled = pointer.From(props.IsReconnectCloneEnabled)
				model.IsRefreshableClone = pointer.From(props.IsRefreshableClone)
				model.TimeUntilReconnectClone = pointer.From(props.TimeUntilReconnectCloneEnabled)

				if props.RefreshableModel != nil {
					model.RefreshableModel = string(*props.RefreshableModel)
				}
				if props.RefreshableStatus != nil {
					model.RefreshableStatus = string(*props.RefreshableStatus)
				}

				// Base properties
				model.AutonomousDatabaseId = pointer.From(props.AutonomousDatabaseId)
				model.BackupRetentionPeriodInDays = pointer.From(props.BackupRetentionPeriodInDays)
				model.CharacterSet = pointer.From(props.CharacterSet)
				model.ComputeCount = pointer.From(props.ComputeCount)
				model.ComputeModel = string(pointer.From(props.ComputeModel))
				model.CustomerContacts = flattenAdbsCustomerContacts(props.CustomerContacts)
				model.DataStorageSizeInGbs = pointer.From(props.DataStorageSizeInGbs)
				model.DataStorageSizeInTbs = pointer.From(props.DataStorageSizeInTbs)
				model.DbVersion = pointer.From(props.DbVersion)
				model.DbWorkload = string(pointer.From(props.DbWorkload))
				model.DisplayName = pointer.From(props.DisplayName)
				model.LicenseModel = string(pointer.From(props.LicenseModel))
				model.AutoScalingEnabled = pointer.From(props.IsAutoScalingEnabled)
				model.AutoScalingForStorageEnabled = pointer.From(props.IsAutoScalingForStorageEnabled)
				model.MtlsConnectionRequired = pointer.From(props.IsMtlsConnectionRequired)
				model.NationalCharacterSet = pointer.From(props.NcharacterSet)
				model.SubnetId = pointer.From(props.SubnetId)
				model.VnetId = pointer.From(props.VnetId)
				model.PrivateEndpoint = pointer.From(props.PrivateEndpoint)
				model.PrivateEndpointIp = pointer.From(props.PrivateEndpointIP)
				model.ServiceConsoleUrl = pointer.From(props.ServiceConsoleURL)
				model.SqlWebDeveloperUrl = pointer.From(props.SqlWebDeveloperURL)
				model.TimeCreated = pointer.From(props.TimeCreated)
				model.OciUrl = pointer.From(props.OciURL)

				if props.LifecycleState != nil {
					model.LifecycleState = string(*props.LifecycleState)
				}

				// Extract connection strings if available
				if props.ConnectionStrings != nil && props.ConnectionStrings.AllConnectionStrings != nil {
					connStrings := []string{}
					allConnStrings := *props.ConnectionStrings.AllConnectionStrings
					if allConnStrings.High != nil {
						connStrings = append(connStrings, *allConnStrings.High)
					}
					if allConnStrings.Medium != nil {
						connStrings = append(connStrings, *allConnStrings.Medium)
					}
					if allConnStrings.Low != nil {
						connStrings = append(connStrings, *allConnStrings.Low)
					}
					model.ConnectionStrings = connStrings
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&model)
		},
	}
}
