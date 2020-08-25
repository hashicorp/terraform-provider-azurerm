package datalake

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataLakeStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDateLakeStoreCreate,
		Read:   resourceArmDateLakeStoreRead,
		Update: resourceArmDateLakeStoreUpdate,
		Delete: resourceArmDateLakeStoreDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataLakeAccountName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"tier": {
				Type:             schema.TypeString,
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
				Type:     schema.TypeString,
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(account.ServiceManaged),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"firewall_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(account.FirewallStateEnabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(account.FirewallStateEnabled),
					string(account.FirewallStateDisabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"firewall_allow_azure_ips": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(account.FirewallAllowAzureIpsStateEnabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(account.FirewallAllowAzureIpsStateEnabled),
					string(account.FirewallAllowAzureIpsStateDisabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDateLakeStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Lake Store %q (Resource Group %q): %s", name, resourceGroup, err)
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

	log.Printf("[INFO] preparing arguments for Data Lake Store creation %q (Resource Group %q)", name, resourceGroup)

	dateLakeStore := account.CreateDataLakeStoreAccountParameters{
		Location: &location,
		Tags:     tags.Expand(t),
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

	future, err := client.Create(ctx, resourceGroup, name, dateLakeStore)
	if err != nil {
		return fmt.Errorf("Error issuing create request for Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Data Lake Store %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmDateLakeStoreRead(d, meta)
}

func resourceArmDateLakeStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
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

	future, err := client.Update(ctx, resourceGroup, name, props)
	if err != nil {
		return fmt.Errorf("Error issuing update request for Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the update of Data Lake Store %q (Resource Group %q) to commplete: %+v", name, resourceGroup, err)
	}

	return resourceArmDateLakeStoreRead(d, meta)
}

func resourceArmDateLakeStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["accounts"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Data Lake Store Account %q was not found (Resource Group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
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

func resourceArmDateLakeStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.StoreAccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["accounts"]
	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error issuing delete request for Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Data Lake Store %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}
