package datalake

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
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
				Default:          string(account.Consumption),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(account.Consumption),
					string(account.Commitment1TB),
					string(account.Commitment10TB),
					string(account.Commitment100TB),
					string(account.Commitment500TB),
					string(account.Commitment1PB),
					string(account.Commitment5PB),
				}, true),
			},

			"encryption_state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(account.Enabled),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(account.Enabled),
					string(account.Disabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"encryption_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(account.ServiceManaged),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"firewall_state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(account.FirewallStateEnabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(account.FirewallStateEnabled),
					string(account.FirewallStateDisabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"firewall_allow_azure_ips": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(account.FirewallAllowAzureIpsStateEnabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(account.FirewallAllowAzureIpsStateEnabled),
					string(account.FirewallAllowAzureIpsStateDisabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, false),
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDateLakeStoreCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	subscriptionId := meta.(*clients.Client).Datalake.StoreAccountsClient.SubscriptionID
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Data Lake Store %s: %+v", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_lake_store", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	tier := d.Get("tier").(string)

	encryptionState := account.EncryptionState(d.Get("encryption_state").(string))
	encryptionType := account.EncryptionConfigType(d.Get("encryption_type").(string))
	firewallState := account.FirewallState(d.Get("firewall_state").(string))
	firewallAllowAzureIPs := account.FirewallAllowAzureIpsState(d.Get("firewall_allow_azure_ips").(string))
	t := d.Get("tags").(map[string]interface{})

	log.Printf("[INFO] preparing arguments for Data Lake Store creation %s", id)

	dateLakeStore := account.CreateDataLakeStoreAccountParameters{
		Location: &location,
		Tags:     tags.Expand(t),
		Identity: expandDataLakeStoreIdentity(d.Get("identity").([]interface{})),
		CreateDataLakeStoreAccountProperties: &account.CreateDataLakeStoreAccountProperties{
			NewTier:               account.TierType(tier),
			FirewallState:         firewallState,
			FirewallAllowAzureIps: firewallAllowAzureIPs,
			EncryptionState:       encryptionState,

			EncryptionConfig: &account.EncryptionConfig{
				Type: encryptionType,
			},
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.Name, dateLakeStore)
	if err != nil {
		return fmt.Errorf("issuing create request for Data Lake Store %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("creating Data Lake Store %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDateLakeStoreRead(d, meta)
}

func resourceArmDateLakeStoreUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	tier := d.Get("tier").(string)
	firewallState := account.FirewallState(d.Get("firewall_state").(string))
	firewallAllowAzureIPs := account.FirewallAllowAzureIpsState(d.Get("firewall_allow_azure_ips").(string))
	t := d.Get("tags").(map[string]interface{})

	props := account.UpdateDataLakeStoreAccountParameters{
		UpdateDataLakeStoreAccountProperties: &account.UpdateDataLakeStoreAccountProperties{
			NewTier:               account.TierType(tier),
			FirewallState:         firewallState,
			FirewallAllowAzureIps: firewallAllowAzureIPs,
		},
		Tags: tags.Expand(t),
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, props)
	if err != nil {
		return fmt.Errorf("issuing update request for Data Lake Store %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of Data Lake Store %s to complete: %+v", id, err)
	}

	return resourceArmDateLakeStoreRead(d, meta)
}

func resourceArmDateLakeStoreRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Data Lake Store Account %s was not found", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Azure Data Lake Store %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenDataLakeStoreIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("flattening identity on Data Lake Store %s: %+v", id, err)
	}

	if properties := resp.DataLakeStoreAccountProperties; properties != nil {
		d.Set("tier", string(properties.CurrentTier))

		d.Set("encryption_state", string(properties.EncryptionState))
		d.Set("firewall_state", string(properties.FirewallState))
		d.Set("firewall_allow_azure_ips", string(properties.FirewallAllowAzureIps))

		if config := properties.EncryptionConfig; config != nil {
			d.Set("encryption_type", string(config.Type))
		}

		d.Set("endpoint", properties.Endpoint)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDateLakeStoreDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccountID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Data Lake Store %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Data Lake Store %s: %+v", id, err)
	}

	return nil
}

func expandDataLakeStoreIdentity(input []interface{}) *account.EncryptionIdentity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &account.EncryptionIdentity{
		Type: utils.String(v["type"].(string)),
	}
}

func flattenDataLakeStoreIdentity(identity *account.EncryptionIdentity) []interface{} {
	if identity == nil {
		return []interface{}{}
	}

	principalID := ""
	if identity.PrincipalID != nil {
		principalID = identity.PrincipalID.String()
	}

	tenantID := ""
	if identity.TenantID != nil {
		tenantID = identity.TenantID.String()
	}

	return []interface{}{
		map[string]interface{}{
			"type":         identity.Type,
			"principal_id": principalID,
			"tenant_id":    tenantID,
		},
	}
}
