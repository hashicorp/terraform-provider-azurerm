// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStreamAnalyticsOutputSql() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsOutputSqlCreateUpdate,
		Read:   resourceStreamAnalyticsOutputSqlRead,
		Update: resourceStreamAnalyticsOutputSqlCreateUpdate,
		Delete: resourceStreamAnalyticsOutputSqlDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := outputs.ParseOutputID(id)
			return err
		}, importStreamAnalyticsOutput(outputs.AzureSqlDatabaseOutputDataSource{})),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsOutputMsSqlV0ToV1{},
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

			"resource_group_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"server": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"database": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"table": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"user": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"max_batch_count": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      10000,
				ValidateFunc: validation.FloatBetween(1, 1073741824),
			},

			"max_writer_count": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.FloatBetween(0, 1),
			},

			"authentication_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(outputs.AuthenticationModeConnectionString),
				ValidateFunc: validation.StringInSlice([]string{
					string(outputs.AuthenticationModeMsi),
					string(outputs.AuthenticationModeConnectionString),
				}, false),
			},
		},
	}
}

func resourceStreamAnalyticsOutputSqlCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := outputs.NewOutputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_stream_analytics_output_mssql", id.ID())
		}
	}

	dataSourceProperties := outputs.AzureSqlDatabaseDataSourceProperties{
		Server:             utils.String(d.Get("server").(string)),
		Database:           utils.String(d.Get("database").(string)),
		Table:              utils.String(d.Get("table").(string)),
		MaxBatchCount:      utils.Float(d.Get("max_batch_count").(float64)),
		MaxWriterCount:     utils.Float(d.Get("max_writer_count").(float64)),
		AuthenticationMode: utils.ToPtr(outputs.AuthenticationMode(d.Get("authentication_mode").(string))),
	}

	// Add user/password dataSourceProperties only if authentication mode requires them
	if *dataSourceProperties.AuthenticationMode == outputs.AuthenticationModeConnectionString {
		dataSourceProperties.User = utils.String(d.Get("user").(string))
		dataSourceProperties.Password = utils.String(d.Get("password").(string))
	}

	props := outputs.Output{
		Name: utils.String(id.OutputName),
		Properties: &outputs.OutputProperties{
			Datasource: &outputs.AzureSqlDatabaseOutputDataSource{
				Properties: &dataSourceProperties,
			},
		},
	}

	var createOpts outputs.CreateOrReplaceOperationOptions
	var updateOpts outputs.UpdateOperationOptions
	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, id, props, createOpts); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
	} else if _, err := client.Update(ctx, id, props, updateOpts); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsOutputSqlRead(d, meta)
}

func resourceStreamAnalyticsOutputSqlRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outputs.ParseOutputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retreving %s: %+v", id, err)
	}

	d.Set("name", id.OutputName)
	d.Set("stream_analytics_job_name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			output, ok := props.Datasource.(outputs.AzureSqlDatabaseOutputDataSource)
			if !ok {
				return fmt.Errorf("converting %s to a SQL Output", *id)
			}

			server := ""
			if v := output.Properties.Server; v != nil {
				server = *v
			}
			d.Set("server", server)

			database := ""
			if v := output.Properties.Database; v != nil {
				database = *v
			}
			d.Set("database", database)

			table := ""
			if v := output.Properties.Table; v != nil {
				table = *v
			}
			d.Set("table", table)

			user := ""
			if v := output.Properties.User; v != nil {
				user = *v
			}
			d.Set("user", user)

			authMode := ""
			if v := output.Properties.AuthenticationMode; v != nil {
				authMode = string(*v)
			}
			d.Set("authentication_mode", authMode)

			maxBatchCount := float64(10000)
			if v := output.Properties.MaxBatchCount; v != nil {
				maxBatchCount = *v
			}
			d.Set("max_batch_count", maxBatchCount)

			maxWriterCount := float64(1)
			if v := output.Properties.MaxWriterCount; v != nil {
				maxWriterCount = *v
			}
			d.Set("max_writer_count", maxWriterCount)
		}
	}
	return nil
}

func resourceStreamAnalyticsOutputSqlDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outputs.ParseOutputID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}
