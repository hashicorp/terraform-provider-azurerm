package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-06-01/netapp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	aznetapp "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetAppAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetAppAccountCreateUpdate,
		Read:   resourceArmNetAppAccountRead,
		Update: resourceArmNetAppAccountCreateUpdate,
		Delete: resourceArmNetAppAccountDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: aznetapp.ValidateNetAppAccountName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocation(),

			"active_directory": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_servers": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: aznetapp.ValidateActiveDirectoryDNSName,
							},
						},
						"domain": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: aznetapp.ValidateActiveDirectoryDomainName,
						},
						"smb_server_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"username": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validate.NoEmptyStrings,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								maskedNew := strings.Replace(new, new, "****************", -1)
								return (new == d.Get(k).(string)) && (maskedNew == old)
							},
						},
						"organizational_unit": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetAppAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Netapp.AccountClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing NetApp Account %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_netapp_account", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	activeDirectories := d.Get("active_directory").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	accountParameters := netapp.Account{
		Location: utils.String(location),
		AccountProperties: &netapp.AccountProperties{
			ActiveDirectories: expandArmNetAppActiveDirectories(activeDirectories),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, accountParameters, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error creating NetApp Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of NetApp Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving NetApp Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read NetApp Account %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmNetAppAccountRead(d, meta)
}

func resourceArmNetAppAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Netapp.AccountClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["netAppAccounts"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NetApp Accounts %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NetApp Accounts %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.AccountProperties; props != nil {
		if err := d.Set("active_directory", flattenArmNetAppActiveDirectories(props.ActiveDirectories)); err != nil {
			return fmt.Errorf("Error setting `active_directory`: %+v", err)
		}
	}

	// Handles tags being interface{} until https://github.com/Azure/azure-rest-api-specs/issues/7447 is fixed
	currentTags := make(map[string]*string)
	if v := resp.Tags; v != nil {
		tagMap := v.(map[string]interface{})

		for k, v := range tagMap {
			currentTags[k] = utils.String(v.(string))
		}
	}

	return tags.FlattenAndSet(d, currentTags)
}

func resourceArmNetAppAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Netapp.AccountClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["netAppAccounts"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting NetApp Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting NetApp Account %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmNetAppActiveDirectories(input []interface{}) *[]netapp.ActiveDirectory {
	results := make([]netapp.ActiveDirectory, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		dns := strings.Join(*utils.ExpandStringSlice(v["dns_servers"].([]interface{})), ",")

		result := netapp.ActiveDirectory{
			DNS:                utils.String(dns),
			Domain:             utils.String(v["domain"].(string)),
			OrganizationalUnit: utils.String(v["organizational_unit"].(string)),
			Password:           utils.String(v["password"].(string)),
			SmbServerName:      utils.String(v["smb_server_name"].(string)),
			Username:           utils.String(v["username"].(string)),
		}

		results = append(results, result)
	}
	return &results
}

func flattenArmNetAppActiveDirectories(input *[]netapp.ActiveDirectory) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		b := make(map[string]interface{})

		if v := item.DNS; v != nil {
			dnsServers := strings.Split(*v, ",")
			b["dns_servers"] = utils.FlattenStringSlice(&dnsServers)
		}
		if v := item.Domain; v != nil {
			b["domain"] = *v
		}
		if v := item.SmbServerName; v != nil {
			b["smb_server_name"] = *v
		}
		if v := item.Username; v != nil {
			b["username"] = *v
		}
		if v := item.Password; v != nil {
			b["password"] = *v
		}
		if v := item.OrganizationalUnit; v != nil {
			b["organizational_unit"] = *v
		}

		results = append(results, b)
	}

	return results
}
