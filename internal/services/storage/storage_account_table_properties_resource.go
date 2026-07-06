// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/shim"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name storage_account_table_properties -service-package-name storage -compare-values "subscription_id:storage_account_id,resource_group_name:storage_account_id,storage_account_name:storage_account_id" -test-name "corsOnly"

type AccountTablePropertiesResource struct{}

var (
	_ sdk.ResourceWithUpdate               = AccountTablePropertiesResource{}
	_ sdk.ResourceWithIdentityTypeOverride = AccountTablePropertiesResource{}
)

type AccountTablePropertiesModel struct {
	StorageAccountId string                                `json:"storage_account_id" tfschema:"storage_account_id"`
	CorsRule         []AccountTablePropertiesCorsRule      `tfschema:"cors_rule"`
	HourMetrics      []AccountTablePropertiesHourMetrics   `tfschema:"hour_metrics"`
	MinuteMetrics    []AccountTablePropertiesMinuteMetrics `tfschema:"minute_metrics"`
	Logging          []AccountTablePropertiesLogging       `tfschema:"logging"`
}

type AccountTablePropertiesCorsRule struct {
	AllowedOrigins []string `tfschema:"allowed_origins"`
	AllowedMethods []string `tfschema:"allowed_methods"`
	AllowedHeaders []string `tfschema:"allowed_headers"`
	ExposedHeaders []string `tfschema:"exposed_headers"`
	MaxAgeSeconds  int64    `tfschema:"max_age_in_seconds"`
}

type AccountTablePropertiesHourMetrics struct {
	Version             string `tfschema:"version"`
	IncludeAPIS         bool   `tfschema:"include_apis"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

type AccountTablePropertiesMinuteMetrics struct {
	Version             string `tfschema:"version"`
	IncludeAPIS         bool   `tfschema:"include_apis"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

type AccountTablePropertiesLogging struct {
	Version             string `tfschema:"version"`
	Delete              bool   `tfschema:"delete"`
	Read                bool   `tfschema:"read"`
	Write               bool   `tfschema:"write"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

var defaultTableCorsProperties = tables.Cors{
	CorsRule: []tables.CorsRule{},
}

var defaultTableHourMetricsProperties = tables.MetricsConfig{
	Version: "1.0",
	Enabled: false,
	RetentionPolicy: tables.RetentionPolicy{
		Enabled: false,
	},
}

var defaultTableMinuteMetricsProperties = tables.MetricsConfig{
	Version: "1.0",
	Enabled: false,
	RetentionPolicy: tables.RetentionPolicy{
		Enabled: false,
	},
}

var defaultTableLoggingProperties = tables.LoggingConfig{
	Version: "1.0",
	Delete:  false,
	Read:    false,
	Write:   false,
	RetentionPolicy: tables.RetentionPolicy{
		Enabled: false,
	},
}

func (s AccountTablePropertiesResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"cors_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 5,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_origins": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"exposed_headers": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 64,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"allowed_headers": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 64,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"allowed_methods": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"DELETE",
								"GET",
								"HEAD",
								"MERGE",
								"POST",
								"OPTIONS",
								"PUT",
							}, false),
						},
					},
					"max_age_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 2000000000),
					},
				},
			},
			AtLeastOneOf: []string{"minute_metrics", "hour_metrics", "logging", "cors_rule"},
		},

		"hour_metrics": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"include_apis": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"retention_policy_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 365),
					},
				},
			},
			AtLeastOneOf: []string{"minute_metrics", "hour_metrics", "logging", "cors_rule"},
		},
		"logging": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"delete": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"read": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"write": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"retention_policy_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 365),
					},
				},
			},
			AtLeastOneOf: []string{"minute_metrics", "hour_metrics", "logging", "cors_rule"},
		},
		"minute_metrics": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"include_apis": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"retention_policy_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 365),
					},
				},
			},
			AtLeastOneOf: []string{"minute_metrics", "hour_metrics", "logging", "cors_rule"},
		},
	}
}

func (s AccountTablePropertiesResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s AccountTablePropertiesResource) ModelObject() interface{} {
	return &AccountTablePropertiesModel{}
}

func (s AccountTablePropertiesResource) ResourceType() string {
	return "azurerm_storage_account_table_properties"
}

func (s AccountTablePropertiesResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateStorageAccountID
}

func (s AccountTablePropertiesResource) Identity() resourceids.ResourceId {
	return &commonids.StorageAccountId{}
}

func (s AccountTablePropertiesResource) IdentityType() pluginsdk.ResourceTypeForIdentity {
	return pluginsdk.ResourceTypeForIdentityVirtual
}

func (s AccountTablePropertiesResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage
			var model AccountTablePropertiesModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountID, err := commonids.ParseStorageAccountID(model.StorageAccountId)
			if err != nil {
				return err
			}

			// Get the target account to ensure it supports tables
			account, err := storageClient.ResourceManager.StorageAccounts.GetProperties(ctx, *accountID, storageaccounts.DefaultGetPropertiesOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *accountID, err)
			}
			if account.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *accountID)
			}

			if account.Model.Sku == nil || account.Model.Sku.Tier == nil || string(account.Model.Sku.Name) == "" {
				return fmt.Errorf("could not read SKU details for %s", *accountID)
			}

			accountTier := *account.Model.Sku.Tier
			accountReplicationTypeParts := strings.Split(string(account.Model.Sku.Name), "_")
			if len(accountReplicationTypeParts) != 2 {
				return fmt.Errorf("could not read SKU replication type for %s", *accountID)
			}
			accountReplicationType := accountReplicationTypeParts[1]

			accountDetails, err := storageClient.GetAccount(ctx, *accountID)
			if err != nil {
				return err
			}
			if accountDetails == nil {
				return fmt.Errorf("unable to locate %s", *accountID)
			}

			supportLevel := availableFunctionalityForAccount(accountDetails.Kind, accountTier, accountReplicationType)

			if !supportLevel.supportTable {
				return fmt.Errorf("account %s does not support tables", *accountID)
			}

			client, err := storageClient.TablesDataPlaneClient(ctx, *accountDetails, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("creating Tables Data Plane Client for %s: %+v", accountID, err)
			}

			props := DefaultValueForAccountTableProperties()
			expandTablePropertiesModel(&model, &props)

			if err = client.UpdateServiceProperties(ctx, props); err != nil {
				return fmt.Errorf("updating Table Properties for %s: %+v", accountID, err)
			}

			// Poll until properties are confirmed set
			if err = pollForTableProperties(ctx, client, props); err != nil {
				return fmt.Errorf("waiting for Table Properties to be set for %s: %+v", accountID, err)
			}

			metadata.SetID(accountID)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, accountID, pluginsdk.ResourceTypeForIdentityVirtual)
		},
	}
}

func (s AccountTablePropertiesResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage

			var state AccountTablePropertiesModel

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			state.StorageAccountId = id.ID()

			account, err := storageClient.GetAccount(ctx, *id)
			if err != nil {
				return metadata.MarkAsGone(id)
			}
			if account == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			client, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Tables Client for %s: %v", *id, err)
			}

			props, err := client.GetServiceProperties(ctx)
			if err != nil {
				return fmt.Errorf("retrieving Table Properties for %s: %+v", *id, err)
			}

			if props != nil {
				if props.Cors != nil && !reflect.DeepEqual(*props.Cors, &defaultTableCorsProperties) {
					corsRules := make([]AccountTablePropertiesCorsRule, 0)
					for _, rule := range props.Cors.CorsRule {
						corsRule := AccountTablePropertiesCorsRule{
							AllowedOrigins: strings.Split(rule.AllowedOrigins, ","),
							AllowedMethods: strings.Split(rule.AllowedMethods, ","),
							AllowedHeaders: strings.Split(rule.AllowedHeaders, ","),
							ExposedHeaders: strings.Split(rule.ExposedHeaders, ","),
							MaxAgeSeconds:  int64(rule.MaxAgeInSeconds),
						}
						corsRules = append(corsRules, corsRule)
					}
					state.CorsRule = corsRules
				}

				if props.HourMetrics != nil && !reflect.DeepEqual(*props.HourMetrics, defaultTableHourMetricsProperties) {
					state.HourMetrics = []AccountTablePropertiesHourMetrics{
						{
							Version:             props.HourMetrics.Version,
							IncludeAPIS:         pointer.From(props.HourMetrics.IncludeAPIs),
							RetentionPolicyDays: int64(props.HourMetrics.RetentionPolicy.Days),
						},
					}
				}

				if props.MinuteMetrics != nil && !reflect.DeepEqual(*props.MinuteMetrics, defaultTableMinuteMetricsProperties) {
					state.MinuteMetrics = []AccountTablePropertiesMinuteMetrics{
						{
							Version:             props.MinuteMetrics.Version,
							IncludeAPIS:         pointer.From(props.MinuteMetrics.IncludeAPIs),
							RetentionPolicyDays: int64(props.MinuteMetrics.RetentionPolicy.Days),
						},
					}
				}

				if props.Logging != nil && !reflect.DeepEqual(*props.Logging, defaultTableLoggingProperties) {
					state.Logging = []AccountTablePropertiesLogging{
						{
							Version:             props.Logging.Version,
							Delete:              props.Logging.Delete,
							Read:                props.Logging.Read,
							Write:               props.Logging.Write,
							RetentionPolicyDays: int64(props.Logging.RetentionPolicy.Days),
						},
					}
				}
			}

			if err = pluginsdk.SetResourceIdentityData(metadata.ResourceData, id, pluginsdk.ResourceTypeForIdentityVirtual); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func (s AccountTablePropertiesResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			account, err := storageClient.GetAccount(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if account == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			client, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Tables Client for %s: %v", *id, err)
			}

			if err = client.UpdateServiceProperties(ctx, DefaultValueForAccountTableProperties()); err != nil {
				return fmt.Errorf("updating Table Properties for %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (s AccountTablePropertiesResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			account, err := storageClient.GetAccount(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if account == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			client, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Tables Client for %s: %v", *id, err)
			}

			props, err := client.GetServiceProperties(ctx)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var model AccountTablePropertiesModel

			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("cors_rule") {
				if len(model.CorsRule) >= 1 {
					corsRules := make([]tables.CorsRule, 0)
					for _, corsRule := range model.CorsRule {
						corsRules = append(corsRules, tables.CorsRule{
							AllowedOrigins:  strings.Join(corsRule.AllowedOrigins, ","),
							AllowedMethods:  strings.Join(corsRule.AllowedMethods, ","),
							AllowedHeaders:  strings.Join(corsRule.AllowedHeaders, ","),
							ExposedHeaders:  strings.Join(corsRule.ExposedHeaders, ","),
							MaxAgeInSeconds: int(corsRule.MaxAgeSeconds),
						})
					}

					if props.Cors == nil {
						props.Cors = &tables.Cors{}
					}
					props.Cors.CorsRule = corsRules
				} else {
					props.Cors = pointer.To(defaultTableCorsProperties)
				}
			}

			if metadata.ResourceData.HasChange("hour_metrics") {
				if len(model.HourMetrics) == 1 {
					metrics := model.HourMetrics[0]
					if props.HourMetrics == nil {
						props.HourMetrics = pointer.To(defaultTableHourMetricsProperties)
					}
					props.HourMetrics.Enabled = true
					if metadata.ResourceData.HasChange("hour_metrics.0.version") {
						props.HourMetrics.Version = metrics.Version
					}

					if metadata.ResourceData.HasChange("hour_metrics.0.include_apis") {
						props.HourMetrics.IncludeAPIs = pointer.To(metrics.IncludeAPIS)
					}

					if metadata.ResourceData.HasChange("hour_metrics.0.retention_policy_days") {
						props.HourMetrics.RetentionPolicy = tables.RetentionPolicy{
							Days:    int(metrics.RetentionPolicyDays),
							Enabled: true,
						}
					}
				} else {
					props.HourMetrics = pointer.To(defaultTableHourMetricsProperties)
				}
			}

			if metadata.ResourceData.HasChange("minute_metrics") {
				if len(model.MinuteMetrics) == 1 {
					metrics := model.MinuteMetrics[0]
					if props.MinuteMetrics == nil {
						props.MinuteMetrics = pointer.To(defaultTableMinuteMetricsProperties)
					}
					props.MinuteMetrics.Enabled = true
					if metadata.ResourceData.HasChange("minute_metrics.0.version") {
						props.MinuteMetrics.Version = metrics.Version
					}

					if metadata.ResourceData.HasChange("minute_metrics.0.include_apis") {
						props.MinuteMetrics.IncludeAPIs = pointer.To(metrics.IncludeAPIS)
					}

					if metadata.ResourceData.HasChange("minute_metrics.0.retention_policy_days") {
						props.MinuteMetrics.RetentionPolicy = tables.RetentionPolicy{
							Days:    int(metrics.RetentionPolicyDays),
							Enabled: true,
						}
					}
				} else {
					props.MinuteMetrics = pointer.To(defaultTableMinuteMetricsProperties)
				}
			}

			if metadata.ResourceData.HasChange("logging") {
				if len(model.Logging) == 1 {
					logging := model.Logging[0]
					if props.Logging == nil {
						props.Logging = pointer.To(defaultTableLoggingProperties)
					}
					if metadata.ResourceData.HasChange("logging.0.version") {
						props.Logging.Version = logging.Version
					}
					if metadata.ResourceData.HasChange("logging.0.delete") {
						props.Logging.Delete = logging.Delete
					}
					if metadata.ResourceData.HasChange("logging.0.read") {
						props.Logging.Read = logging.Read
					}
					if metadata.ResourceData.HasChange("logging.0.write") {
						props.Logging.Write = logging.Write
					}
					if metadata.ResourceData.HasChange("logging.0.retention_policy_days") {
						props.Logging.RetentionPolicy = tables.RetentionPolicy{
							Days:    int(logging.RetentionPolicyDays),
							Enabled: true,
						}
					}
				} else {
					props.Logging = pointer.To(defaultTableLoggingProperties)
				}
			}

			if err = client.UpdateServiceProperties(ctx, *props); err != nil {
				return fmt.Errorf("updating Table Properties for %s: %+v", *id, err)
			}

			// Poll until properties are confirmed set
			if err = pollForTableProperties(ctx, client, *props); err != nil {
				return fmt.Errorf("waiting for Table Properties to be set for %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandTablePropertiesModel(model *AccountTablePropertiesModel, props *tables.StorageServiceProperties) {
	if len(model.CorsRule) >= 1 {
		corsRules := make([]tables.CorsRule, 0)
		for _, corsRule := range model.CorsRule {
			corsRules = append(corsRules, tables.CorsRule{
				AllowedOrigins:  strings.Join(corsRule.AllowedOrigins, ","),
				AllowedMethods:  strings.Join(corsRule.AllowedMethods, ","),
				AllowedHeaders:  strings.Join(corsRule.AllowedHeaders, ","),
				ExposedHeaders:  strings.Join(corsRule.ExposedHeaders, ","),
				MaxAgeInSeconds: int(corsRule.MaxAgeSeconds),
			})
		}

		props.Cors.CorsRule = corsRules
	}

	if len(model.HourMetrics) == 1 {
		metrics := model.HourMetrics[0]
		props.HourMetrics.Enabled = true
		props.HourMetrics.Version = metrics.Version
		if metrics.RetentionPolicyDays != 0 {
			props.HourMetrics.RetentionPolicy = tables.RetentionPolicy{
				Days:    int(metrics.RetentionPolicyDays),
				Enabled: true,
			}
		}

		props.HourMetrics.IncludeAPIs = pointer.To(metrics.IncludeAPIS)
	}

	if len(model.MinuteMetrics) != 0 {
		metrics := model.MinuteMetrics[0]
		props.MinuteMetrics.Enabled = true
		props.MinuteMetrics.Version = metrics.Version
		if metrics.RetentionPolicyDays != 0 {
			props.MinuteMetrics.RetentionPolicy = tables.RetentionPolicy{
				Days:    int(metrics.RetentionPolicyDays),
				Enabled: true,
			}
		}

		props.MinuteMetrics.IncludeAPIs = pointer.To(metrics.IncludeAPIS)
	}

	if len(model.Logging) != 0 {
		logging := model.Logging[0]
		props.Logging.Version = logging.Version
		props.Logging.Delete = logging.Delete
		props.Logging.Read = logging.Read
		props.Logging.Write = logging.Write
		if logging.RetentionPolicyDays != 0 {
			props.Logging.RetentionPolicy = tables.RetentionPolicy{
				Enabled: true,
				Days:    int(logging.RetentionPolicyDays),
			}
		}
	}
}

func pollForTableProperties(ctx context.Context, client shim.StorageTableWrapper, expected tables.StorageServiceProperties) error {
	pollerType := custompollers.NewDataPlaneTablesPropertiesPoller(client, expected)
	propertiesPoller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	return propertiesPoller.PollUntilDone(ctx)
}

func DefaultValueForAccountTableProperties() tables.StorageServiceProperties {
	return tables.StorageServiceProperties{
		Logging: &tables.LoggingConfig{
			Version: "1.0",
			Delete:  false,
			Read:    false,
			Write:   false,
			RetentionPolicy: tables.RetentionPolicy{
				Enabled: false,
			},
		},
		HourMetrics: &tables.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: tables.RetentionPolicy{
				Enabled: false,
			},
		},
		MinuteMetrics: &tables.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: tables.RetentionPolicy{
				Enabled: false,
			},
		},
		Cors: &tables.Cors{
			CorsRule: []tables.CorsRule{},
		},
	}
}
