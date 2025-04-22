// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/elasticpools"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlElasticPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlElasticPoolCreateUpdate,
		Read:   resourceMsSqlElasticPoolRead,
		Update: resourceMsSqlElasticPoolCreateUpdate,
		Delete: resourceMsSqlElasticPoolDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ElasticPoolID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlElasticPoolName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlServerName,
			},

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"BasicPool",
								"StandardPool",
								"PremiumPool",
								"GP_Gen4",
								"GP_Gen5",
								"GP_Fsv2",
								"GP_DC",
								"BC_Gen4",
								"BC_Gen5",
								"BC_DC",
								"HS_Gen5",
								"HS_PRMS",
								"HS_MOPRMS",
							}, false),
						},

						"capacity": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"tier": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Basic",
								"Standard",
								"Premium",
								"GeneralPurpose",
								"BusinessCritical",
								"Hyperscale",
							}, false),
						},

						"family": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen4",
								"Gen5",
								"Fsv2",
								"DC",
								"MOPRMS",
								"PRMS",
							}, false),
						},
					},
				},
			},

			"maintenance_configuration_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "SQL_Default",
				ValidateFunc: validation.StringInSlice(resourceMsSqlDatabaseMaintenanceNames(), false),
			},

			"per_database_settings": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"min_capacity": {
							Type:         pluginsdk.TypeFloat,
							Required:     true,
							ValidateFunc: validation.FloatAtLeast(0.0),
						},

						"max_capacity": {
							Type:         pluginsdk.TypeFloat,
							Required:     true,
							ValidateFunc: validation.FloatAtLeast(0.0),
						},
					},
				},
			},

			"max_size_bytes": {
				Type:          pluginsdk.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"max_size_gb"},
				ValidateFunc:  validation.IntAtLeast(0),
			},

			"max_size_gb": {
				Type:          pluginsdk.TypeFloat,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"max_size_bytes"},
				ValidateFunc:  validation.FloatAtLeast(0),
			},

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			// NOTE: The implementation of 'enclave_type' in the API differs slightly between database
			// and elasticpools. Database does not allow the 'Default' value to be passed for DW or
			// DC skus, where elasticpools allows 'Default' but will error if you try to set the
			// 'enclave_type' to 'VBS' for DC skus...
			"enclave_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true, // TODO: Remove Computed in 4.0
				ValidateFunc: validation.StringInSlice([]string{
					string(databases.AlwaysEncryptedEnclaveTypeVBS),
					string(databases.AlwaysEncryptedEnclaveTypeDefault),
				}, false),
			},

			"license_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(elasticpools.ElasticPoolLicenseTypeBasePrice),
					string(elasticpools.ElasticPoolLicenseTypeLicenseIncluded),
				}, false),
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
				if err := helper.MSSQLElasticPoolValidateSKU(diff); err != nil {
					return err
				}

				return nil
			},

			pluginsdk.ForceNewIfChange("enclave_type", func(ctx context.Context, old, new, _ interface{}) bool {
				// enclave_type cannot be removed once it has been set
				// but can be changed between VBS and Default...
				// this Diff will not work until 4.0 when we remove
				// the computed property from the field scheam.
				if old.(string) != "" && new.(string) == "" {
					return true
				}

				return false
			}),
		),
	}
}

func resourceMsSqlElasticPoolCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MSSQL ElasticPool creation.")

	id := commonids.NewSqlElasticPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mssql_elasticpool", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := expandMsSqlElasticPoolSku(d)

	maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(subscriptionId, d.Get("maintenance_configuration_name").(string))
	elasticPool := elasticpools.ElasticPool{
		Name:     pointer.To(id.ElasticPoolName),
		Location: location,
		Sku:      sku,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &elasticpools.ElasticPoolProperties{
			LicenseType:                pointer.To(elasticpools.ElasticPoolLicenseType(d.Get("license_type").(string))),
			PerDatabaseSettings:        expandMsSqlElasticPoolPerDatabaseSettings(d),
			ZoneRedundant:              pointer.To(d.Get("zone_redundant").(bool)),
			MaintenanceConfigurationId: pointer.To(maintenanceConfigId.ID()),
			PreferredEnclaveType:       nil,
		},
	}

	// NOTE: The service default is actually nil/empty which indicates enclave is disabled. the value `Default` is NOT the default.
	if v, ok := d.GetOk("enclave_type"); ok && v.(string) != "" {
		elasticPool.Properties.PreferredEnclaveType = pointer.To(elasticpools.AlwaysEncryptedEnclaveType(v.(string)))
	}

	if d.HasChange("max_size_gb") {
		if v, ok := d.GetOk("max_size_gb"); ok {
			maxSizeBytes := v.(float64) * 1073741824
			elasticPool.Properties.MaxSizeBytes = utils.Int64(int64(maxSizeBytes))
		}
	} else if v, ok := d.GetOk("max_size_bytes"); ok {
		elasticPool.Properties.MaxSizeBytes = pointer.To(int64(v.(int)))
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, elasticPool)
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceMsSqlElasticPoolRead(d, meta)
}

func resourceMsSqlElasticPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSqlElasticPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("name", model.Name)
		d.Set("resource_group_name", id.ResourceGroupName)
		d.Set("location", model.Location)
		d.Set("server_name", id.ServerName)

		if err := d.Set("sku", flattenMsSqlElasticPoolSku(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if props := model.Properties; props != nil {
			enclaveType := ""
			if v := props.PreferredEnclaveType; v != nil {
				enclaveType = string(pointer.From(v))
			}
			d.Set("enclave_type", enclaveType)

			// Basic tier does not return max_size_bytes, so we need to skip setting this
			// value if the pricing tier is equal to Basic
			if tier, ok := d.GetOk("sku.0.tier"); ok {
				if !strings.EqualFold(tier.(string), "Basic") {
					d.Set("max_size_gb", pointer.To(*props.MaxSizeBytes/int64(1073741824)))
					d.Set("max_size_bytes", pointer.To(props.MaxSizeBytes))
				}
			}

			d.Set("zone_redundant", pointer.From(props.ZoneRedundant))

			licenseType := string(elasticpools.ElasticPoolLicenseTypeLicenseIncluded)
			if props.LicenseType != nil {
				licenseType = string(*props.LicenseType)
			}
			d.Set("license_type", licenseType)

			if err := d.Set("per_database_settings", flattenMsSqlElasticPoolPerDatabaseSettings(props.PerDatabaseSettings)); err != nil {
				return fmt.Errorf("setting `per_database_settings`: %+v", err)
			}

			maintenanceConfigId, err := publicmaintenanceconfigurations.ParsePublicMaintenanceConfigurationIDInsensitively(*props.MaintenanceConfigurationId)
			if err != nil {
				return err
			}
			d.Set("maintenance_configuration_name", maintenanceConfigId.PublicMaintenanceConfigurationName)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceMsSqlElasticPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSqlElasticPoolID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting ElasticPool %s: %+v", id, err)
	}

	return nil
}

func expandMsSqlElasticPoolPerDatabaseSettings(d *pluginsdk.ResourceData) *elasticpools.ElasticPoolPerDatabaseSettings {
	perDatabaseSettings := d.Get("per_database_settings").([]interface{})
	perDatabaseSetting := perDatabaseSettings[0].(map[string]interface{})

	minCapacity := perDatabaseSetting["min_capacity"].(float64)
	maxCapacity := perDatabaseSetting["max_capacity"].(float64)

	return &elasticpools.ElasticPoolPerDatabaseSettings{
		MinCapacity: utils.Float(minCapacity),
		MaxCapacity: utils.Float(maxCapacity),
	}
}

func expandMsSqlElasticPoolSku(d *pluginsdk.ResourceData) *elasticpools.Sku {
	skus := d.Get("sku").([]interface{})
	sku := skus[0].(map[string]interface{})

	name := sku["name"].(string)
	tier := sku["tier"].(string)
	family := sku["family"].(string)
	capacity := sku["capacity"].(int)

	return &elasticpools.Sku{
		Name:     name,
		Tier:     pointer.To(tier),
		Family:   pointer.To(family),
		Capacity: pointer.To(int64(capacity)),
	}
}

func flattenMsSqlElasticPoolSku(input *elasticpools.Sku) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	values := map[string]interface{}{}

	if name := input.Name; name != "" {
		values["name"] = name
	}

	if tier := input.Tier; tier != nil {
		values["tier"] = *tier
	}

	if family := input.Family; family != nil {
		values["family"] = *family
	}

	if capacity := input.Capacity; capacity != nil {
		values["capacity"] = *capacity
	}

	return []interface{}{values}
}

func flattenMsSqlElasticPoolPerDatabaseSettings(resp *elasticpools.ElasticPoolPerDatabaseSettings) []interface{} {
	perDatabaseSettings := map[string]interface{}{}

	if minCapacity := resp.MinCapacity; minCapacity != nil {
		perDatabaseSettings["min_capacity"] = *minCapacity
	}

	if maxCapacity := resp.MaxCapacity; maxCapacity != nil {
		perDatabaseSettings["max_capacity"] = *maxCapacity
	}

	return []interface{}{perDatabaseSettings}
}
