package netapp

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-10-01/netapp"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetAppAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetAppAccountCreateUpdate,
		Read:   resourceArmNetAppAccountRead,
		Update: resourceArmNetAppAccountCreateUpdate,
		Delete: resourceArmNetAppAccountDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.NetAppAccountID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateNetAppAccountName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

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
								ValidateFunc: validate.IPv4Address,
							},
						},
						"domain": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^[(\da-zA-Z-).]{1,255}$`),
								`The domain name must end with a letter or number before dot and start with a letter or number after dot and can not be longer than 255 characters in length.`,
							),
						},
						"smb_server_name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^[\da-zA-Z]{1,10}$`),
								`The smb server name can not be longer than 10 characters in length.`,
							),
						},
						"username": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"organizational_unit": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetAppAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
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

	accountParameters := netapp.Account{
		Location: utils.String(location),
		AccountProperties: &netapp.AccountProperties{
			ActiveDirectories: expandArmNetAppActiveDirectories(activeDirectories),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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
	client := meta.(*clients.Client).NetApp.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NetApp Accounts %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NetApp Accounts %q (Resource Group %q): %+v", id.NetAppAccountName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmNetAppAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.NetAppAccountName)
	if err != nil {
		return fmt.Errorf("Error deleting NetApp Account %q (Resource Group %q): %+v", id.NetAppAccountName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting NetApp Account %q (Resource Group %q): %+v", id.NetAppAccountName, id.ResourceGroup, err)
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
