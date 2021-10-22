package datalake

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datalake/analytics/mgmt/2016-11-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
				Default:          string(account.Consumption),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(account.Consumption),
					string(account.Commitment100000AUHours),
					string(account.Commitment10000AUHours),
					string(account.Commitment1000AUHours),
					string(account.Commitment100AUHours),
					string(account.Commitment500000AUHours),
					string(account.Commitment50000AUHours),
					string(account.Commitment5000AUHours),
					string(account.Commitment500AUHours),
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
	subscriptionId := meta.(*clients.Client).Datalake.AnalyticsAccountsClient.SubscriptionID
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAnalyticsAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.AccountName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Data Lake Analytics Account %s: %+v", id, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_data_lake_analytics_account", *existing.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	storeAccountName := d.Get("default_store_account_name").(string)
	tier := d.Get("tier").(string)
	t := d.Get("tags").(map[string]interface{})

	log.Printf("[INFO] preparing arguments for Azure ARM Date Lake Store creation %s", id)

	dateLakeAnalyticsAccount := account.CreateDataLakeAnalyticsAccountParameters{
		Location: &location,
		Tags:     tags.Expand(t),
		CreateDataLakeAnalyticsAccountProperties: &account.CreateDataLakeAnalyticsAccountProperties{
			NewTier:                     account.TierType(tier),
			DefaultDataLakeStoreAccount: &storeAccountName,
			DataLakeStoreAccounts: &[]account.AddDataLakeStoreWithAccountParameters{
				{
					Name: &storeAccountName,
				},
			},
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.AccountName, dateLakeAnalyticsAccount)
	if err != nil {
		return fmt.Errorf("issuing create request for Data Lake Analytics Account %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("creating Data Lake Analytics Account %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDateLakeAnalyticsAccountRead(d, meta)
}

func resourceArmDateLakeAnalyticsAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.AnalyticsAccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AnalyticsAccountID(d.Id())
	if err != nil {
		return err
	}

	storeAccountName := d.Get("default_store_account_name").(string)
	newTier := d.Get("tier").(string)
	newTags := d.Get("tags").(map[string]interface{})

	props := &account.UpdateDataLakeAnalyticsAccountParameters{
		Tags: tags.Expand(newTags),
		UpdateDataLakeAnalyticsAccountProperties: &account.UpdateDataLakeAnalyticsAccountProperties{
			NewTier: account.TierType(newTier),
			DataLakeStoreAccounts: &[]account.UpdateDataLakeStoreWithAccountParameters{
				{
					Name: &storeAccountName,
				},
			},
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.AccountName, props)
	if err != nil {
		return fmt.Errorf("issuing update request for Data Lake Analytics Account %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of Data Lake Analytics Account %s to complete: %+v", id, err)
	}

	return resourceArmDateLakeAnalyticsAccountRead(d, meta)
}

func resourceArmDateLakeAnalyticsAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.AnalyticsAccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AnalyticsAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] DataLakeAnalyticsAccountAccount '%s'", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Azure Data Lake Analytics Account %s: %+v", id, err)
	}

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if properties := resp.DataLakeAnalyticsAccountProperties; properties != nil {
		d.Set("tier", string(properties.CurrentTier))
		d.Set("default_store_account_name", properties.DefaultDataLakeStoreAccount)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDateLakeAnalyticsAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.AnalyticsAccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AnalyticsAccountID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.AccountName)
	if err != nil {
		return fmt.Errorf("deleting Data Lake Analytics Account %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of Data Lake Analytics Account %s: %+v", id, err)
	}

	return nil
}
