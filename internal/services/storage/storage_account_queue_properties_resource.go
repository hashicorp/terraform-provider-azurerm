// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
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

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:CREATE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	model := existing.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	// TODO: Add Import error support

	accountTier := pointer.From(model.Sku.Tier)
	accountKind := pointer.From(model.Kind)
	replicationType := strings.ToUpper(strings.Split(string(model.Sku.Name), "_")[1])
	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, replicationType)

	if !supportLevel.supportQueue {
		return fmt.Errorf("%q are not supported for account kind %q in sku tier %q", storageAccountQueuePropertiesResourceName, accountKind, accountTier)
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	dataPlaneAccount, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if dataPlaneAccount == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	// NOTE: Wait for the data plane queue container to become available...
	log.Printf("[DEBUG] [%s:CREATE] Calling 'custompollers.NewDataPlaneQueuesAvailabilityPoller' building Queues Poller: %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	pollerType, err := custompollers.NewDataPlaneQueuesAvailabilityPoller(ctx, storageClient, dataPlaneAccount)
	if err != nil {
		return fmt.Errorf("building Queues Poller: %+v", err)
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'poller.PollUntilDone' waiting for the Queues Service to become available: %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for the Queues Service to become available: %+v", err)
	}

	// NOTE: Now that we know the data plane container is available, we can now set the properties on the resource...
	log.Printf("[DEBUG] [%s:CREATE] Calling 'storageClient.QueuesDataPlaneClient' building Queues Client: %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	queuesDataPlaneClient, err := storageClient.QueuesDataPlaneClient(ctx, *dataPlaneAccount, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Queues Client: %s", err)
	}

	queueProperties, err := standaloneExpandAccountQueueProperties(d)
	if err != nil {
		return fmt.Errorf("expanding `properties`: %+v", err)
	}

	log.Printf("[DEBUG] [%s:CREATE] Calling 'queuesDataPlaneClient.UpdateServiceProperties': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	if err = queuesDataPlaneClient.UpdateServiceProperties(ctx, *queueProperties); err != nil {
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

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:UPDATE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := existing.Model
	if err := validateExistingModel(model, id); err != nil {
		return err
	}

	accountKind := pointer.From(model.Kind)
	accountReplicationType := strings.ToUpper(strings.Split(string(model.Sku.Name), "_")[1])
	accountTier := pointer.From(model.Sku.Tier)
	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)

	if d.HasChange("cors_rule") || d.HasChange("logging") || d.HasChange("minute_metrics") || d.HasChange("hour_metrics") {
		if !supportLevel.supportQueue {
			return fmt.Errorf("%q are not supported for a storage account with the account kind %q in sku tier %q", storageAccountQueuePropertiesResourceName, accountKind, accountTier)
		}

		log.Printf("[DEBUG] [%s:UPDATE] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
		account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if account == nil {
			return fmt.Errorf("unable to locate %s", *id)
		}

		log.Printf("[DEBUG] [%s:UPDATE] Calling 'storageClient.QueuesDataPlaneClient': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
		queueDataPlaneClient, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Queues Client: %s", err)
		}

		queueProperties, err := standaloneExpandAccountQueueProperties(d)
		if err != nil {
			return fmt.Errorf("expanding `properties` for %s: %+v", *id, err)
		}

		log.Printf("[DEBUG] [%s:UPDATE] Calling 'queueClient.UpdateServiceProperties': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
		if err = queueDataPlaneClient.UpdateServiceProperties(ctx, *queueProperties); err != nil {
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

	log.Printf("[DEBUG] [%s:READ] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	log.Printf("[DEBUG] [%s:READ] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	log.Printf("[DEBUG] [%s:READ] Calling 'storageClient.QueuesDataPlaneClient': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	queueDataPlaneClient, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Queues Client: %s", err)
	}

	log.Printf("[DEBUG] [%s:READ] Calling 'queueDataPlaneClient.GetServiceProperties': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	queueProps, err := queueDataPlaneClient.GetServiceProperties(ctx)
	if err != nil {
		return fmt.Errorf("retrieving queue properties for %s: %+v", *id, err)
	}

	if err := d.Set("cors_rule", flattenAccountQueuePropertiesCors(queueProps.Cors)); err != nil {
		return fmt.Errorf("setting `cors_rule`: %+v", err)
	}

	if queueProps.Logging != nil {
		if err := d.Set("logging", standaloneFlattenAccountQueuePropertiesLogging(queueProps.Logging)); err != nil {
			return fmt.Errorf("setting `logging`: %+v", err)
		}
	}

	if queueProps.HourMetrics != nil {
		if err := d.Set("hour_metrics", standaloneFlattenAccountQueuePropertiesMetrics(queueProps.HourMetrics)); err != nil {
			return fmt.Errorf("setting `hour_metrics`: %+v", err)
		}
	}

	if queueProps.MinuteMetrics != nil {
		if err := d.Set("minute_metrics", standaloneFlattenAccountQueuePropertiesMetrics(queueProps.MinuteMetrics)); err != nil {
			return fmt.Errorf("setting `minute_metrics`: %+v", err)
		}
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

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	log.Printf("[DEBUG] [%s:DELETE] Calling 'client.GetProperties': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	_, err = client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	log.Printf("[DEBUG] [%s:DELETE] Calling 'storageClient.FindAccount': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	dataPlaneAccount, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if dataPlaneAccount == nil {
		return fmt.Errorf("unable to locate %s", *id)
	}

	log.Printf("[DEBUG] [%s:DELETE] Calling 'storageClient.QueuesDataPlaneClient': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)
	queueDataPlaneClient, err := storageClient.QueuesDataPlaneClient(ctx, *dataPlaneAccount, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Plane Queues Client: %s", err)
	}

	log.Printf("[DEBUG] [%s:DELETE] Calling 'queueDataPlaneClient.UpdateServiceProperties': %s", strings.ToUpper(storageAccountQueuePropertiesResourceName), id)

	// NOTE: Since this is a fake resource that has been split off from the main storage account resource
	// the best we can do is reset the values to the default settings...
	queueProperties := defaultAccountQueueProperties()
	if err = queueDataPlaneClient.UpdateServiceProperties(ctx, queueProperties); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func standaloneExpandAccountQueueProperties(d *pluginsdk.ResourceData) (*queues.StorageServiceProperties, error) {
	var err error
	properties := defaultAccountQueueProperties()

	corsRule := d.Get("cors_rule").([]interface{})
	logging := d.Get("logging").([]interface{})
	minuteMetrics := d.Get("minute_metrics").([]interface{})
	hourMetrics := d.Get("hour_metrics").([]interface{})

	if len(corsRule) != 0 {
		properties.Cors = expandAccountQueuePropertiesCors(corsRule)
	}

	if len(logging) != 0 {
		if v := logging[0].(map[string]interface{}); v != nil {
			properties.Logging = expandAccountQueuePropertiesLogging([]interface{}{v})
		}
	}

	if len(minuteMetrics) != 0 {
		if v := minuteMetrics[0].(map[string]interface{}); v != nil {
			properties.MinuteMetrics, err = expandAccountQueuePropertiesMetrics([]interface{}{v})
			if err != nil {
				return nil, fmt.Errorf("expanding `minute_metrics`: %+v", err)
			}
		}
	}

	if len(hourMetrics) != 0 {
		if v := hourMetrics[0].(map[string]interface{}); v != nil {
			properties.HourMetrics, err = expandAccountQueuePropertiesMetrics([]interface{}{v})
			if err != nil {
				return nil, fmt.Errorf("expanding `hour_metrics`: %+v", err)
			}
		}
	}

	return &properties, nil
}

// TODO: Remove in v5.0, this is only here for legacy support of existing Storage Accounts...
func expandAccountQueueProperties(input []interface{}) (*queues.StorageServiceProperties, error) {
	var err error
	properties := defaultAccountQueueProperties()

	if len(input) != 0 {
		attrs := input[0].(map[string]interface{})

		if attrs["cors_rule"] != nil {
			properties.Cors = expandAccountQueuePropertiesCors(attrs["cors_rule"].([]interface{}))
		}

		if attrs["logging"] != nil {
			properties.Logging = expandAccountQueuePropertiesLogging(attrs["logging"].([]interface{}))
		}

		if attrs["minute_metrics"] != nil {
			properties.MinuteMetrics, err = expandAccountQueuePropertiesMetrics(attrs["minute_metrics"].([]interface{}))
			if err != nil {
				return nil, fmt.Errorf("expanding `minute_metrics`: %+v", err)
			}
		}

		if attrs["hour_metrics"] != nil {
			properties.HourMetrics, err = expandAccountQueuePropertiesMetrics(attrs["hour_metrics"].([]interface{}))
			if err != nil {
				return nil, fmt.Errorf("expanding `hour_metrics`: %+v", err)
			}
		}
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
				"logging":        logging,
				"hour_metrics":   hourMetrics,
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

func standaloneFlattenAccountQueuePropertiesLogging(input *queues.LoggingConfig) []interface{} {
	output := []interface{}{}

	if input == nil || (input.Version == "1.0" && !input.Delete && !input.Read && !input.Write && input.RetentionPolicy.Days == 0) {
		return output
	}

	retentionPolicyDays := input.RetentionPolicy.Days

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

// TODO: Remove in v5.0, this is only here for legacy support of existing Storage Accounts...
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

func standaloneFlattenAccountQueuePropertiesMetrics(input *queues.MetricsConfig) []interface{} {
	output := make([]interface{}, 0)

	if input == nil || (input.Version == "1.0" && !input.Enabled && !input.RetentionPolicy.Enabled && input.RetentionPolicy.Days == 0 && input.IncludeAPIs == nil) {
		return output
	}

	if input.Version != "" {
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

// TODO: Remove in v5.0, this is only here for legacy support of existing Storage Accounts...
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

func defaultAccountQueueProperties() queues.StorageServiceProperties {
	output := queues.StorageServiceProperties{
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

	return output
}
