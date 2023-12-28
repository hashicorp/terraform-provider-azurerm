// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsReferenceMsSql() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsReferenceInputMsSqlCreateUpdate,
		Read:   resourceStreamAnalyticsReferenceInputMsSqlRead,
		Update: resourceStreamAnalyticsReferenceInputMsSqlCreateUpdate,
		Delete: resourceStreamAnalyticsReferenceInputMsSqlDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := inputs.ParseInputID(id)
			return err
		}, importStreamAnalyticsReferenceInput("Microsoft.Sql/Server/Database")),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsReferenceInputMsSqlV0ToV1{},
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"stream_analytics_job_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"server": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"database": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"username": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"password": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"refresh_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Static",
					"RefreshPeriodicallyWithFull",
					"RefreshPeriodicallyWithDelta",
				}, false),
			},

			"refresh_interval_duration": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.BatchMaxWaitTime,
			},

			"full_snapshot_query": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"delta_snapshot_query": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"table": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceStreamAnalyticsReferenceInputMsSqlCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Reference Input MsSql creation.")
	id := inputs.NewInputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_stream_analytics_reference_input_mssql", id.ID())
		}
	}

	refreshType := d.Get("refresh_type").(string)

	if _, ok := d.GetOk("refresh_interval_duration"); refreshType != "Static" && !ok {
		return fmt.Errorf("refresh_interval_duration must be set if refresh_type is RefreshPeriodicallyWithFull or RefreshPeriodicallyWithDelta")
	} else if _, ok = d.GetOk("delta_snapshot_query"); refreshType == "Static" && ok {
		return fmt.Errorf("delta_snapshot_query cannot be set if refresh_type is Static")
	}

	properties := &inputs.AzureSqlReferenceInputDataSourceProperties{
		Server:      pointer.To(d.Get("server").(string)),
		Database:    pointer.To(d.Get("database").(string)),
		User:        pointer.To(d.Get("username").(string)),
		Password:    pointer.To(d.Get("password").(string)),
		RefreshType: pointer.To(inputs.RefreshType(refreshType)),
	}

	if v, ok := d.GetOk("refresh_interval_duration"); ok {
		properties.RefreshRate = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("full_snapshot_query"); ok {
		properties.FullSnapshotQuery = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("delta_snapshot_query"); ok {
		properties.DeltaSnapshotQuery = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("table"); ok {
		properties.Table = pointer.To(v.(string))
	}

	props := inputs.Input{
		Name: pointer.To(id.InputName),
		Properties: &inputs.ReferenceInputProperties{
			Datasource: &inputs.AzureSqlReferenceInputDataSource{
				Properties: properties,
			},
		},
	}

	var opts inputs.CreateOrReplaceOperationOptions
	if _, err := client.CreateOrReplace(ctx, id, props, opts); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStreamAnalyticsReferenceInputMsSqlRead(d, meta)
}

func resourceStreamAnalyticsReferenceInputMsSqlRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.InputName)
	d.Set("stream_analytics_job_name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			input, ok := props.(inputs.InputProperties) // nolint: gosimple
			if !ok {
				return fmt.Errorf("converting %s to an Input", *id)
			}

			reference, ok := input.(inputs.ReferenceInputProperties)
			if !ok {
				return fmt.Errorf("converting %s to Reference Input", *id)
			}

			referenceInputAzureSql, ok := reference.Datasource.(inputs.AzureSqlReferenceInputDataSource)
			if !ok {
				return fmt.Errorf("converting %s to an Azure Sql Reference Input", *id)
			}

			if referenceInputAzureSql.Properties != nil {
				server := ""
				if v := referenceInputAzureSql.Properties.Server; v != nil {
					server = *v
				}
				d.Set("server", server)

				database := ""
				if v := referenceInputAzureSql.Properties.Database; v != nil {
					database = *v
				}
				d.Set("database", database)

				username := ""
				if v := referenceInputAzureSql.Properties.User; v != nil {
					username = *v
				}
				d.Set("username", username)

				refreshType := ""
				if v := referenceInputAzureSql.Properties.RefreshType; v != nil {
					refreshType = string(*v)
				}
				d.Set("refresh_type", refreshType)

				intervalDuration := ""
				if v := referenceInputAzureSql.Properties.RefreshRate; v != nil {
					intervalDuration = *v
				}
				d.Set("refresh_interval_duration", intervalDuration)

				fullSnapshotQuery := ""
				if v := referenceInputAzureSql.Properties.FullSnapshotQuery; v != nil {
					fullSnapshotQuery = *v
				}
				d.Set("full_snapshot_query", fullSnapshotQuery)

				deltaSnapshotQuery := ""
				if v := referenceInputAzureSql.Properties.DeltaSnapshotQuery; v != nil {
					deltaSnapshotQuery = *v
				}
				d.Set("delta_snapshot_query", deltaSnapshotQuery)

				table := ""
				if v := referenceInputAzureSql.Properties.Table; v != nil {
					table = *v
				}
				d.Set("table", table)
			}
		}
	}

	return nil
}

func resourceStreamAnalyticsReferenceInputMsSqlDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}
	return nil
}
