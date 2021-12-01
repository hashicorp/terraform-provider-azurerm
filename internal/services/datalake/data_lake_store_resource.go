package datalake

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/sdk/datalakestore/2016-11-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataLakeStore() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmDateLakeStoreCreate,
		Read:   resourceArmDateLakeStoreRead,
		Update: resourceArmDateLakeStoreUpdate,
		Delete: resourceArmDateLakeStoreDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AccountID(id)
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
					string(accounts.TierTypeCommitmentOneTB),
					string(accounts.TierTypeCommitmentOneZeroTB),
					string(accounts.TierTypeCommitmentOneZeroZeroTB),
					string(accounts.TierTypeCommitmentFiveZeroZeroTB),
					string(accounts.TierTypeCommitmentOnePB),
					string(accounts.TierTypeCommitmentFivePB),
				}, true),
			},

			"encryption_state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(accounts.EncryptionStateEnabled),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(accounts.EncryptionStateEnabled),
					string(accounts.EncryptionStateDisabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"encryption_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(accounts.EncryptionConfigTypeServiceManaged),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"firewall_state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(accounts.FirewallStateEnabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(accounts.FirewallStateEnabled),
					string(accounts.FirewallStateDisabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"firewall_allow_azure_ips": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(accounts.FirewallAllowAzureIpsStateEnabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(accounts.FirewallAllowAzureIpsStateEnabled),
					string(accounts.FirewallAllowAzureIpsStateDisabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedIdentity(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDateLakeStoreCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	subscriptionId := meta.(*clients.Client).Datalake.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := accounts.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_lake_store", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	tier := accounts.TierType(d.Get("tier").(string))

	encryptionState := accounts.EncryptionState(d.Get("encryption_state").(string))
	encryptionType := accounts.EncryptionConfigType(d.Get("encryption_type").(string))
	firewallState := accounts.FirewallState(d.Get("firewall_state").(string))
	firewallAllowAzureIPs := accounts.FirewallAllowAzureIpsState(d.Get("firewall_allow_azure_ips").(string))
	t := d.Get("tags").(map[string]interface{})

	log.Printf("[INFO] preparing arguments for Data Lake Store creation %s", id)

	identity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	dateLakeStore := accounts.CreateDataLakeStoreAccountParameters{
		Location: location,
		Tags:     tagsHelper.Expand(t),
		Identity: identity,
		Properties: &accounts.CreateDataLakeStoreAccountProperties{
			NewTier:               &tier,
			FirewallState:         &firewallState,
			FirewallAllowAzureIps: &firewallAllowAzureIPs,
			EncryptionState:       &encryptionState,

			EncryptionConfig: &accounts.EncryptionConfig{
				Type: encryptionType,
			},
		},
	}

	if err := client.CreateThenPoll(ctx, id, dateLakeStore); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceArmDateLakeStoreRead(d, meta)
}

func resourceArmDateLakeStoreUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	tier := accounts.TierType(d.Get("tier").(string))
	firewallState := accounts.FirewallState(d.Get("firewall_state").(string))
	firewallAllowAzureIPs := accounts.FirewallAllowAzureIpsState(d.Get("firewall_allow_azure_ips").(string))
	t := d.Get("tags").(map[string]interface{})

	props := accounts.UpdateDataLakeStoreAccountParameters{
		Properties: &accounts.UpdateDataLakeStoreAccountProperties{
			NewTier:               &tier,
			FirewallState:         &firewallState,
			FirewallAllowAzureIps: &firewallAllowAzureIPs,
		},
		Tags: tagsHelper.Expand(t),
	}

	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceArmDateLakeStoreRead(d, meta)
}

func resourceArmDateLakeStoreRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] Data Lake Store Account %s was not found", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retreiving %s: %+v", id, err)
	}

	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if err := d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("flattening identity on Data Lake Store %s: %+v", id, err)
		}

		if properties := model.Properties; properties != nil {
			tier := ""
			if properties.CurrentTier != nil {
				tier = string(*properties.CurrentTier)
			}
			d.Set("tier", tier)

			encryptionState := ""
			if properties.EncryptionState != nil {
				encryptionState = string(*properties.EncryptionState)
			}
			d.Set("encryption_state", encryptionState)

			firewallState := ""
			if properties.FirewallState != nil {
				firewallState = string(*properties.FirewallState)
			}
			d.Set("firewall_state", firewallState)

			firewallAllowAzureIps := ""
			if properties.FirewallAllowAzureIps != nil {
				firewallAllowAzureIps = string(*properties.FirewallAllowAzureIps)
			}
			d.Set("firewall_allow_azure_ips", firewallAllowAzureIps)

			if config := properties.EncryptionConfig; config != nil {
				d.Set("encryption_type", string(config.Type))
			}

			d.Set("endpoint", properties.Endpoint)
		}

		return tags.FlattenAndSet(d, tagsHelper.Flatten(model.Tags))
	}
	return nil
}

func resourceArmDateLakeStoreDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
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
