// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/queue/queues"
)

var storageAccountQueuePropertiesResourceName = "azurerm_storage_account_queue_properties"

func resourceStorageAccountQueueProperties() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageAccountQueuePropertiesCreate,
		Read:   resourceStorageAccountQueuePropertiesRead,
		Update: resourceStorageAccountQueuePropertiesUpdate,
		Delete: resourceStorageAccountQueuePropertiesDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseStorageAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"properties": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cors_rule": helpers.SchemaStorageAccountCorsRule(false),

						"hour_metrics": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"version": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									// TODO 4.0: Remove this property and determine whether to enable based on existence of the out side block.
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},
									"include_apis": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},
									"retention_policy_days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
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
									// TODO 4.0: Remove this property and determine whether to enable based on existence of the out side block.
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},
									"include_apis": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},
									"retention_policy_days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceStorageAccountQueuePropertiesCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountQueuePropertiesResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountQueuePropertiesResourceName)

	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	// NOTE: Import error cannot be supported for this resource...

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Kind == nil {
		return fmt.Errorf("retrieving %s: `model.Kind` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}

	if existing.Model.Sku == nil {
		return fmt.Errorf("retrieving %s: `model.Sku` was nil", id)
	}

	var accountKind storageaccounts.Kind
	var accountTier storageaccounts.SkuTier
	accountReplicationType := ""

	accountKind = *existing.Model.Kind
	accountReplicationType = strings.Split(string(existing.Model.Sku.Name), "_")[1]
	if existing.Model.Sku.Tier != nil {
		accountTier = *existing.Model.Sku.Tier
	}

	dataPlaneAccount, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if dataPlaneAccount == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)

	if !supportLevel.supportQueue {
		return fmt.Errorf("`properties` are not supported for account kind %q in sku tier %q", accountKind, accountTier)
	}

	queueClient, err := storageClient.QueuesDataPlaneClient(ctx, *dataPlaneAccount, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Queues Client: %s", err)
	}

	queueProperties, err := expandAccountQueueProperties(d.Get("properties").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `properties`: %+v", err)
	}

	if err = queueClient.UpdateServiceProperties(ctx, *queueProperties); err != nil {
		return fmt.Errorf("creating `properties`: %+v", err)
	}

	d.SetId(id.ID())

	return resourceStorageAccountQueuePropertiesRead(d, meta)
}

func resourceStorageAccountQueuePropertiesUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountQueuePropertiesResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountQueuePropertiesResourceName)

	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Kind == nil {
		return fmt.Errorf("retrieving %s: `model.Kind` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}
	if existing.Model.Sku == nil {
		return fmt.Errorf("retrieving %s: `model.Sku` was nil", id)
	}

	var accountKind storageaccounts.Kind
	var accountTier storageaccounts.SkuTier
	accountReplicationType := ""

	accountKind = *existing.Model.Kind
	accountReplicationType = strings.Split(string(existing.Model.Sku.Name), "_")[1]
	if existing.Model.Sku.Tier != nil {
		accountTier = *existing.Model.Sku.Tier
	}

	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)

	if d.HasChange("properties") {
		if !supportLevel.supportQueue {
			return fmt.Errorf("queue properties are not supported for a storage account with the account kind %q in sku tier %q", accountKind, accountTier)
		}

		account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if account == nil {
			return fmt.Errorf("unable to locate %s", *id)
		}

		queueClient, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Queues Client: %s", err)
		}

		queueProperties, err := expandAccountQueueProperties(d.Get("properties").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `properties` for %s: %+v", *id, err)
		}

		if err = queueClient.UpdateServiceProperties(ctx, *queueProperties); err != nil {
			return fmt.Errorf("updating `properties` for %s: %+v", *id, err)
		}
	}

	return resourceStorageAccountQueuePropertiesRead(d, meta)
}

func resourceStorageAccountQueuePropertiesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	queueClient, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Queues Client: %s", err)
	}

	queueProps, err := queueClient.GetServiceProperties(ctx)
	if err != nil {
		return fmt.Errorf("retrieving queue properties for %s: %+v", *id, err)
	}

	queueProperties := flattenAccountQueueProperties(queueProps)

	if err := d.Set("properties", queueProperties); err != nil {
		return fmt.Errorf("setting `properties`: %+v", err)
	}

	return nil
}

func resourceStorageAccountQueuePropertiesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountQueuePropertiesResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountQueuePropertiesResourceName)

	_, err = client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	dataPlaneAccount, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if dataPlaneAccount == nil {
		return fmt.Errorf("unable to locate %s", *id)
	}

	queueClient, err := storageClient.QueuesDataPlaneClient(ctx, *dataPlaneAccount, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Plane Queues Client: %s", err)
	}

	// NOTE: Call expand with an empty interface to get an
	// unconfigured block back from the function...
	queueProperties, err := expandAccountQueueProperties(make([]interface{}, 0))
	if err != nil {
		return fmt.Errorf("expanding %s: %+v", *id, err)
	}

	if err = queueClient.UpdateServiceProperties(ctx, *queueProperties); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandAccountQueueProperties(input []interface{}) (*queues.StorageServiceProperties, error) {
	var err error
	properties := queues.StorageServiceProperties{
		Cors: &queues.Cors{
			CorsRule: []queues.CorsRule{},
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
		Logging: &queues.LoggingConfig{
			Version: "1.0",
			Delete:  false,
			Read:    false,
			Write:   false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		},
	}

	if len(input) == 0 {
		return &properties, nil
	}

	attrs := input[0].(map[string]interface{})

	properties.Cors = expandAccountQueuePropertiesCors(attrs["cors_rule"].([]interface{}))
	properties.Logging = expandAccountQueuePropertiesLogging(attrs["logging"].([]interface{}))

	properties.MinuteMetrics, err = expandAccountQueuePropertiesMetrics(attrs["minute_metrics"].([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("expanding `minute_metrics`: %+v", err)
	}

	properties.HourMetrics, err = expandAccountQueuePropertiesMetrics(attrs["hour_metrics"].([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("expanding `hour_metrics`: %+v", err)
	}

	return &properties, nil
}

func flattenAccountQueueProperties(input *queues.StorageServiceProperties) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		corsRules := flattenAccountQueuePropertiesCors(input.Cors)
		logging := flattenAccountQueuePropertiesLogging(input.Logging)
		hourMetrics := flattenAccountQueuePropertiesMetrics(input.HourMetrics)
		minuteMetrics := flattenAccountQueuePropertiesMetrics(input.MinuteMetrics)

		if len(corsRules) > 0 || len(logging) > 0 || len(hourMetrics) > 0 || len(minuteMetrics) > 0 {
			output = append(output, map[string]interface{}{
				"cors_rule":      corsRules,
				"hour_metrics":   hourMetrics,
				"logging":        logging,
				"minute_metrics": minuteMetrics,
			})
		}
	}

	return output
}

func expandAccountQueuePropertiesLogging(input []interface{}) *queues.LoggingConfig {
	if len(input) == 0 {
		return &queues.LoggingConfig{
			Version: "1.0",
			Delete:  false,
			Read:    false,
			Write:   false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		}
	}

	loggingAttr := input[0].(map[string]interface{})
	logging := &queues.LoggingConfig{
		Delete:  loggingAttr["delete"].(bool),
		Read:    loggingAttr["read"].(bool),
		Version: loggingAttr["version"].(string),
		Write:   loggingAttr["write"].(bool),
	}

	if v, ok := loggingAttr["retention_policy_days"]; ok {
		if days := v.(int); days > 0 {
			logging.RetentionPolicy = queues.RetentionPolicy{
				Days:    days,
				Enabled: true,
			}
		}
	}

	return logging
}

func flattenAccountQueuePropertiesLogging(input *queues.LoggingConfig) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	retentionPolicyDays := 0
	if input.RetentionPolicy.Enabled {
		retentionPolicyDays = input.RetentionPolicy.Days
	}

	return []interface{}{
		map[string]interface{}{
			"delete":                input.Delete,
			"read":                  input.Read,
			"retention_policy_days": retentionPolicyDays,
			"version":               input.Version,
			"write":                 input.Write,
		},
	}
}

func expandAccountQueuePropertiesMetrics(input []interface{}) (*queues.MetricsConfig, error) {
	if len(input) == 0 {
		return &queues.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		}, nil
	}

	metricsAttr := input[0].(map[string]interface{})

	metrics := &queues.MetricsConfig{
		Enabled: metricsAttr["enabled"].(bool),
		Version: metricsAttr["version"].(string),
	}

	if v, ok := metricsAttr["retention_policy_days"]; ok {
		if days := v.(int); days > 0 {
			metrics.RetentionPolicy = queues.RetentionPolicy{
				Days:    days,
				Enabled: true,
			}
		}
	}

	if v, ok := metricsAttr["include_apis"]; ok {
		includeAPIs := v.(bool)
		if metrics.Enabled {
			metrics.IncludeAPIs = &includeAPIs
		} else if includeAPIs {
			return nil, fmt.Errorf("`include_apis` may only be set when `enabled` is true")
		}
	}

	return metrics, nil
}

func flattenAccountQueuePropertiesMetrics(input *queues.MetricsConfig) []interface{} {
	output := make([]interface{}, 0)

	if input != nil && input.Version != "" {
		retentionPolicyDays := 0
		if input.RetentionPolicy.Enabled {
			retentionPolicyDays = input.RetentionPolicy.Days
		}

		output = append(output, map[string]interface{}{
			"enabled":               input.Enabled,
			"include_apis":          pointer.From(input.IncludeAPIs),
			"retention_policy_days": retentionPolicyDays,
			"version":               input.Version,
		})
	}

	return output
}

func expandAccountQueuePropertiesCors(input []interface{}) *queues.Cors {
	if len(input) == 0 {
		return &queues.Cors{}
	}

	corsRules := make([]queues.CorsRule, 0)
	for _, attr := range input {
		corsRuleAttr := attr.(map[string]interface{})
		corsRule := queues.CorsRule{}

		corsRule.AllowedOrigins = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_origins"].([]interface{})), ",")
		corsRule.ExposedHeaders = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["exposed_headers"].([]interface{})), ",")
		corsRule.AllowedHeaders = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_headers"].([]interface{})), ",")
		corsRule.AllowedMethods = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_methods"].([]interface{})), ",")
		corsRule.MaxAgeInSeconds = corsRuleAttr["max_age_in_seconds"].(int)

		corsRules = append(corsRules, corsRule)
	}

	cors := &queues.Cors{
		CorsRule: corsRules,
	}
	return cors
}

func flattenAccountQueuePropertiesCors(input *queues.Cors) []interface{} {
	output := make([]interface{}, 0)

	if input == nil || len(input.CorsRule) == 0 || input.CorsRule[0].AllowedOrigins == "" {
		return output
	}

	for _, item := range input.CorsRule {
		output = append(output, map[string]interface{}{
			"allowed_headers":    flattenAccountQueuePropertiesCorsRule(item.AllowedHeaders),
			"allowed_methods":    flattenAccountQueuePropertiesCorsRule(item.AllowedMethods),
			"allowed_origins":    flattenAccountQueuePropertiesCorsRule(item.AllowedOrigins),
			"exposed_headers":    flattenAccountQueuePropertiesCorsRule(item.ExposedHeaders),
			"max_age_in_seconds": item.MaxAgeInSeconds,
		})
	}

	return output
}

func flattenAccountQueuePropertiesCorsRule(input string) []interface{} {
	results := make([]interface{}, 0)

	components := strings.Split(input, ",")
	for _, item := range components {
		results = append(results, item)
	}

	return results
}
