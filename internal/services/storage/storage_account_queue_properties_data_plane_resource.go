// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
)

type AccountQueuePropertiesResource struct{}

var _ sdk.ResourceWithUpdate = AccountQueuePropertiesResource{}

type AccountQueuePropertiesModel struct {
	StorageAccountId string                                `json:"storage_account_id" tfschema:"storage_account_id"`
	CorsRule         []AccountQueuePropertiesCorsRule      `tfschema:"cors_rule"`
	HourMetrics      []AccountQueuePropertiesHourMetrics   `tfschema:"hour_metrics"`
	MinuteMetrics    []AccountQueuePropertiesMinuteMetrics `tfschema:"minute_metrics"`
	Logging          []AccountQueuePropertiesLogging       `tfschema:"logging"`
}

type AccountQueuePropertiesCorsRule struct {
	AllowedOrigins []string `tfschema:"allowed_origins"`
	AllowedMethods []string `tfschema:"allowed_methods"`
	AllowedHeaders []string `tfschema:"allowed_headers"`
	ExposedHeaders []string `tfschema:"exposed_headers"`
	MaxAgeSeconds  int64    `tfschema:"max_age_in_seconds"`
}

type AccountQueuePropertiesHourMetrics struct {
	Version             string `tfschema:"version"`
	IncludeAPIS         bool   `tfschema:"include_apis"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

type AccountQueuePropertiesMinuteMetrics struct {
	Version             string `tfschema:"version"`
	IncludeAPIS         bool   `tfschema:"include_apis"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

type AccountQueuePropertiesLogging struct {
	Version             string `tfschema:"version"`
	Delete              bool   `tfschema:"delete"`
	Read                bool   `tfschema:"read"`
	Write               bool   `tfschema:"write"`
	RetentionPolicyDays int64  `tfschema:"retention_policy_days"`
}

var defaultCorsProperties = queues.Cors{
	CorsRule: []queues.CorsRule{},
}

var defaultHourMetricsProperties = queues.MetricsConfig{
	Version: "1.0",
	Enabled: false,
	RetentionPolicy: queues.RetentionPolicy{
		Enabled: false,
	},
}

var defaultMinuteMetricsProperties = queues.MetricsConfig{
	Version: "1.0",
	Enabled: false,
	RetentionPolicy: queues.RetentionPolicy{
		Enabled: false,
	},
}

var defaultLoggingProperties = queues.LoggingConfig{
	Version: "1.0",
	Delete:  false,
	Read:    false,
	Write:   false,
	RetentionPolicy: queues.RetentionPolicy{
		Enabled: false,
	},
}

func (s AccountQueuePropertiesResource) Arguments() map[string]*pluginsdk.Schema {
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

func (s AccountQueuePropertiesResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s AccountQueuePropertiesResource) ModelObject() interface{} {
	return &AccountQueuePropertiesModel{}
}

func (s AccountQueuePropertiesResource) ResourceType() string {
	return "azurerm_storage_account_queue_properties"
}

func (s AccountQueuePropertiesResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateStorageAccountID
}

func (s AccountQueuePropertiesResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage
			var model AccountQueuePropertiesModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountID, err := commonids.ParseStorageAccountID(model.StorageAccountId)
			if err != nil {
				return err
			}

			// Get the target account to ensure it supports queues
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

			accountDetails, err := storageClient.FindAccount(ctx, accountID.SubscriptionId, accountID.StorageAccountName)
			if err != nil {
				return err
			}

			supportLevel := availableFunctionalityForAccount(accountDetails.Kind, accountTier, accountReplicationType)

			if !supportLevel.supportQueue {
				return fmt.Errorf("account %s does not support queues", *accountID)
			}

			client, err := storageClient.QueuesDataPlaneClient(ctx, *accountDetails, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("creating Queues Data Plane Client for %s: %+v", accountID, err)
			}

			props := DefaultValueForAccountQueueProperties()

			if len(model.CorsRule) >= 1 {
				corsRules := make([]queues.CorsRule, 0)
				for _, corsRule := range model.CorsRule {
					corsRules = append(corsRules, queues.CorsRule{
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
					props.HourMetrics.RetentionPolicy = queues.RetentionPolicy{
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
					props.MinuteMetrics.RetentionPolicy = queues.RetentionPolicy{
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
					props.Logging.RetentionPolicy = queues.RetentionPolicy{
						Enabled: true,
						Days:    int(logging.RetentionPolicyDays),
					}
				}
			}

			if err = client.UpdateServiceProperties(ctx, props); err != nil {
				return fmt.Errorf("updating Queue Properties for %s: %+v", accountID, err)
			}

			metadata.SetID(accountID)

			return nil
		},
	}
}

func (s AccountQueuePropertiesResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage

			var state AccountQueuePropertiesModel

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			state.StorageAccountId = id.ID()

			account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
			if err != nil {
				return metadata.MarkAsGone(id)
			}
			if account == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			client, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Queues Client for %s: %v", *id, err)
			}

			props, err := client.GetServiceProperties(ctx)
			if err != nil {
				return fmt.Errorf("retrieving Queue Properties for %s: %+v", *id, err)
			}

			if props != nil {
				if props.Cors != nil && !reflect.DeepEqual(*props.Cors, &defaultCorsProperties) {
					corsRules := make([]AccountQueuePropertiesCorsRule, 0)
					for _, rule := range props.Cors.CorsRule {
						corsRule := AccountQueuePropertiesCorsRule{
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

				if props.HourMetrics != nil && !reflect.DeepEqual(*props.HourMetrics, &defaultHourMetricsProperties) {
					state.HourMetrics = []AccountQueuePropertiesHourMetrics{
						{
							Version:             props.HourMetrics.Version,
							IncludeAPIS:         pointer.From(props.HourMetrics.IncludeAPIs),
							RetentionPolicyDays: int64(props.HourMetrics.RetentionPolicy.Days),
						},
					}
				}

				if props.MinuteMetrics != nil && !reflect.DeepEqual(*props.MinuteMetrics, &defaultMinuteMetricsProperties) {
					state.MinuteMetrics = []AccountQueuePropertiesMinuteMetrics{
						{
							Version:             props.MinuteMetrics.Version,
							IncludeAPIS:         pointer.From(props.MinuteMetrics.IncludeAPIs),
							RetentionPolicyDays: int64(props.MinuteMetrics.RetentionPolicy.Days),
						},
					}
				}

				if props.Logging != nil && !reflect.DeepEqual(*props.Logging, &defaultLoggingProperties) {
					state.Logging = []AccountQueuePropertiesLogging{
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

			return metadata.Encode(&state)
		},
	}
}

func (s AccountQueuePropertiesResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if account == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			client, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Queues Client for %s: %v", *id, err)
			}

			if err = client.UpdateServiceProperties(ctx, DefaultValueForAccountQueueProperties()); err != nil {
				return fmt.Errorf("updating Queue Properties for %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (s AccountQueuePropertiesResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage

			id, err := commonids.ParseStorageAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if account == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			client, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Queues Client for %s: %v", *id, err)
			}

			props, err := client.GetServiceProperties(ctx)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var model AccountQueuePropertiesModel

			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("cors_rule") {
				if len(model.CorsRule) >= 1 {
					corsRules := make([]queues.CorsRule, 0)
					for _, corsRule := range model.CorsRule {
						corsRules = append(corsRules, queues.CorsRule{
							AllowedOrigins:  strings.Join(corsRule.AllowedOrigins, ","),
							AllowedMethods:  strings.Join(corsRule.AllowedMethods, ","),
							AllowedHeaders:  strings.Join(corsRule.AllowedHeaders, ","),
							ExposedHeaders:  strings.Join(corsRule.ExposedHeaders, ","),
							MaxAgeInSeconds: int(corsRule.MaxAgeSeconds),
						})
					}

					props.Cors.CorsRule = corsRules
				} else {
					props.Cors = pointer.To(defaultCorsProperties)
				}
			}

			if metadata.ResourceData.HasChange("hour_metrics") {
				if len(model.HourMetrics) == 1 {
					metrics := model.HourMetrics[0]
					if metadata.ResourceData.HasChange("hour_metrics.0.version") {
						props.HourMetrics.Version = metrics.Version
					}

					if metadata.ResourceData.HasChange("hour_metrics.0.include_apis") {
						props.HourMetrics.IncludeAPIs = pointer.To(metrics.IncludeAPIS)
					}

					if metadata.ResourceData.HasChange("hour_metrics.0.retention_policy_days") {
						props.HourMetrics.RetentionPolicy = queues.RetentionPolicy{
							Days:    int(metrics.RetentionPolicyDays),
							Enabled: true,
						}
					}
				} else {
					props.HourMetrics = pointer.To(defaultHourMetricsProperties)
				}
			}

			if metadata.ResourceData.HasChange("minute_metrics") {
				if len(model.MinuteMetrics) == 1 {
					metrics := model.MinuteMetrics[0]
					if metadata.ResourceData.HasChange("minute_metrics.0.version") {
						props.MinuteMetrics.Version = metrics.Version
					}

					if metadata.ResourceData.HasChange("minute_metrics.0.include_apis") {
						props.MinuteMetrics.IncludeAPIs = pointer.To(metrics.IncludeAPIS)
					}

					if metadata.ResourceData.HasChange("minute_metrics.0.retention_policy_days") {
						props.MinuteMetrics.RetentionPolicy = queues.RetentionPolicy{
							Days:    int(metrics.RetentionPolicyDays),
							Enabled: true,
						}
					}
				} else {
					props.MinuteMetrics = pointer.To(defaultMinuteMetricsProperties)
				}
			}

			if metadata.ResourceData.HasChange("logging") {
				if len(model.Logging) == 1 {
					logging := model.Logging[0]
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
						props.Logging.RetentionPolicy = queues.RetentionPolicy{
							Days:    int(logging.RetentionPolicyDays),
							Enabled: true,
						}
					}
				} else {
					props.Logging = pointer.To(defaultLoggingProperties)
				}
			}

			if err = client.UpdateServiceProperties(ctx, *props); err != nil {
				return fmt.Errorf("updating Queue Properties for %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func DefaultValueForAccountQueueProperties() queues.StorageServiceProperties {
	return queues.StorageServiceProperties{
		Logging: &queues.LoggingConfig{
			Version: "1.0",
			Delete:  false,
			Read:    false,
			Write:   false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		},
		HourMetrics: &queues.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		},
		MinuteMetrics: &queues.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		},
		Cors: &queues.Cors{
			CorsRule: []queues.CorsRule{},
		},
	}
}
