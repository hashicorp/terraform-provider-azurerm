package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-06-01/netapp"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"active_directories": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"domain": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"organizational_unit": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
						"smb_server_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetAppAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.AccountClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing NetApp Account %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_netapp_account", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	activeDirectories := d.Get("active_directories").([]interface{})
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
	client := meta.(*ArmClient).netapp.AccountClient
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
	if accountProperties := resp.AccountProperties; accountProperties != nil {
		if err := d.Set("active_directories", flattenArmNetAppActiveDirectories(accountProperties.ActiveDirectories)); err != nil {
			return fmt.Errorf("Error setting `active_directories`: %+v", err)
		}
	}

	return tags.FlattenAndSetTags(d, resp.Tags)
}

func resourceArmNetAppAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.AccountClient
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
		dns := v["dns"].(string)
		domain := v["domain"].(string)
		organizationalUnit := v["organizational_unit"].(string)
		password := v["password"].(string)
		smbServerName := v["smb_server_name"].(string)
		userName := v["username"].(string)

		result := netapp.ActiveDirectory{
			DNS:                utils.String(dns),
			Domain:             utils.String(domain),
			OrganizationalUnit: utils.String(organizationalUnit),
			Password:           utils.String(password),
			SmbServerName:      utils.String(smbServerName),
			Username:           utils.String(userName),
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
		v := make(map[string]interface{})

		if id := item.ActiveDirectoryID; id != nil {
			v["id"] = *id
		}
		if dns := item.DNS; dns != nil {
			v["dns"] = *dns
		}
		if domain := item.Domain; domain != nil {
			v["domain"] = *domain
		}
		if organizationalUnit := item.OrganizationalUnit; organizationalUnit != nil {
			v["organizational_unit"] = *organizationalUnit
		}
		if smbServerName := item.SmbServerName; smbServerName != nil {
			v["smb_server_name"] = *smbServerName
		}
		if status := item.Status; status != nil {
			v["status"] = *status
		}
		if userName := item.Username; userName != nil {
			v["username"] = *userName
		}
		if password := item.Password; password != nil {
			v["password"] = *password
		}

		results = append(results, v)
	}

	return results
}
