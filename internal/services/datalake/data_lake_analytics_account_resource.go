package datalake

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/sdk/datalakeanalytics/2016-11-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataLakeAnalyticsAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmDateLakeAnalyticsAccountCreate,
		Read:   resourceArmDateLakeAnalyticsAccountRead,
		Update: resourceArmDateLakeAnalyticsAccountUpdate,
		Delete: resourceArmDateLakeAnalyticsAccountDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AnalyticsAccountID(id)
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
				ValidateFunc: validate.AccountName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"tier": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Default:          string(accounts.TierTypeConsumption),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(accounts.TierTypeConsumption),
					string(accounts.TierTypeCommitmentOneZeroZeroZeroZeroZeroAUHours),
					string(accounts.TierTypeCommitmentOneZeroZeroZeroZeroAUHours),
					string(accounts.TierTypeCommitmentOneZeroZeroZeroAUHours),
					string(accounts.TierTypeCommitmentOneZeroZeroAUHours),
					string(accounts.TierTypeCommitmentFiveZeroZeroZeroZeroZeroAUHours),
					string(accounts.TierTypeCommitmentFiveZeroZeroZeroZeroAUHours),
					string(accounts.TierTypeCommitmentFiveZeroZeroZeroAUHours),
					string(accounts.TierTypeCommitmentFiveZeroZeroAUHours),
				}, true),
			},

			"default_store_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName(),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDateLakeAnalyticsAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.AnalyticsAccountsClient
	subscriptionId := meta.(*clients.Client).Datalake.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := accounts.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing Data Lake Analytics Account %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_data_lake_analytics_account", id.ID())
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	storeAccountName := d.Get("default_store_account_name").(string)
	tier := accounts.TierType(d.Get("tier").(string))
	t := d.Get("tags").(map[string]interface{})

	log.Printf("[INFO] preparing arguments for Azure ARM Date Lake Store creation %s", id)

	dateLakeAnalyticsAccount := accounts.CreateDataLakeAnalyticsAccountParameters{
		Location: location,
		Tags:     tagsHelper.Expand(t),
		Properties: accounts.CreateDataLakeAnalyticsAccountProperties{
			NewTier:                     &tier,
			DefaultDataLakeStoreAccount: storeAccountName,
			DataLakeStoreAccounts: []accounts.AddDataLakeStoreWithAccountParameters{
				{
					Name: storeAccountName,
				},
			},
		},
	}

	if err := client.CreateThenPoll(ctx, id, dateLakeAnalyticsAccount); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDateLakeAnalyticsAccountRead(d, meta)
}

func resourceArmDateLakeAnalyticsAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.AnalyticsAccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	storeAccountName := d.Get("default_store_account_name").(string)
	newTier := accounts.TierType(d.Get("tier").(string))
	newTags := d.Get("tags").(map[string]interface{})

	props := accounts.UpdateDataLakeAnalyticsAccountParameters{
		Tags: tagsHelper.Expand(newTags),
		Properties: &accounts.UpdateDataLakeAnalyticsAccountProperties{
			NewTier: &newTier,
			DataLakeStoreAccounts: &[]accounts.UpdateDataLakeStoreWithAccountParameters{
				{
					Name: storeAccountName,
				},
			},
		},
	}

	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceArmDateLakeAnalyticsAccountRead(d, meta)
}

func resourceArmDateLakeAnalyticsAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.AnalyticsAccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] DataLakeAnalyticsAccountAccount '%s'", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if properties := model.Properties; properties != nil {
			tier := ""
			if properties.CurrentTier != nil {
				tier = string(*properties.CurrentTier)
			}
			d.Set("tier", tier)
			d.Set("default_store_account_name", properties.DefaultDataLakeStoreAccount)
		}

		return tags.FlattenAndSet(d, tagsHelper.Flatten(model.Tags))
	}
	return nil
}

func resourceArmDateLakeAnalyticsAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.AnalyticsAccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
