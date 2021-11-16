package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/streamanalytics/mgmt/2020-03-01-preview/streamanalytics"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStreamAnalyticsReferenceMsSql() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsReferenceInputMsSqlCreateUpdate,
		Read:   resourceStreamAnalyticsReferenceInputMsSqlRead,
		Update: resourceStreamAnalyticsReferenceInputMsSqlCreateUpdate,
		Delete: resourceStreamAnalyticsReferenceInputMsSqlDelete,
		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.StreamInputID(id)
			return err
		}, importStreamAnalyticsReferenceInput(streamanalytics.TypeBasicReferenceInputDataSourceTypeMicrosoftSQLServerDatabase)),

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

			"resource_group_name": azure.SchemaResourceGroupName(),

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
		},
	}
}

func resourceStreamAnalyticsReferenceInputMsSqlCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Reference Input MsSql creation.")
	id := parse.NewStreamInputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.InputName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_stream_analytics_reference_input_mssql", id.ID())
		}
	}

	refreshType := d.Get("refresh_type").(string)

	if _, ok := d.GetOk("refresh_interval_duration"); refreshType != "Static" && !ok {
		return fmt.Errorf("refresh_interval_duration must be set if refresh_type is RefreshPeriodicallyWithFull or RefreshPeriodicallyWithDelta")
	} else if _, ok = d.GetOk("delta_snapshot_query"); refreshType == "Static" && ok {
		return fmt.Errorf("delta_snapshot_query cannot be set if refresh_type is Static")
	}

	properties := &streamanalytics.AzureSQLReferenceInputDataSourceProperties{
		Server:      utils.String(d.Get("server").(string)),
		Database:    utils.String(d.Get("database").(string)),
		User:        utils.String(d.Get("username").(string)),
		Password:    utils.String(d.Get("password").(string)),
		RefreshType: utils.String(refreshType),
	}

	if v, ok := d.GetOk("refresh_interval_duration"); ok {
		properties.RefreshRate = utils.String(v.(string))
	}

	if v, ok := d.GetOk("full_snapshot_query"); ok {
		properties.FullSnapshotQuery = utils.String(v.(string))
	}

	if v, ok := d.GetOk("delta_snapshot_query"); ok {
		properties.DeltaSnapshotQuery = utils.String(v.(string))
	}

	props := streamanalytics.Input{
		Name: utils.String(id.InputName),
		Properties: &streamanalytics.ReferenceInputProperties{
			Type: streamanalytics.TypeReference,
			Datasource: &streamanalytics.AzureSQLReferenceInputDataSource{
				Type:       streamanalytics.TypeBasicReferenceInputDataSourceTypeMicrosoftSQLServerDatabase,
				Properties: properties,
			},
		},
	}

	if _, err := client.CreateOrReplace(ctx, props, id.ResourceGroup, id.StreamingjobName, id.InputName, "", ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStreamAnalyticsReferenceInputMsSqlRead(d, meta)
}

func resourceStreamAnalyticsReferenceInputMsSqlRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamInputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.InputName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.InputName)
	d.Set("stream_analytics_job_name", id.StreamingjobName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.Properties; props != nil {
		v, ok := props.AsReferenceInputProperties()
		if !ok {
			return fmt.Errorf("converting Reference Input MS SQL to a Reference Input: %+v", err)
		}

		inputDataSource, ok := v.Datasource.AsAzureSQLReferenceInputDataSource()
		if !ok {
			return fmt.Errorf("converting Reference Input MS SQL to a MS SQL Stream Input: %+v", err)
		}

		d.Set("server", inputDataSource.Properties.Server)
		d.Set("database", inputDataSource.Properties.Database)
		d.Set("username", inputDataSource.Properties.User)
		d.Set("refresh_type", inputDataSource.Properties.RefreshType)
		d.Set("refresh_interval_duration", inputDataSource.Properties.RefreshRate)
		d.Set("full_snapshot_query", inputDataSource.Properties.FullSnapshotQuery)
		d.Set("delta_snapshot_query", inputDataSource.Properties.DeltaSnapshotQuery)
	}
	return nil
}

func resourceStreamAnalyticsReferenceInputMsSqlDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamInputID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.StreamingjobName, id.InputName); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}
	return nil
}
