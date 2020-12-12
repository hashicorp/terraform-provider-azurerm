package compute

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSshPublicKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSshPublicKeyCreate,
		Read:   resourceSshPublicKeyRead,
		Update: resourceSshPublicKeyUpdate,
		Delete: resourceSshPublicKeyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SSHPublicKeyID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Public SSH Key name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"public_key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: ValidateSSHKey,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSshPublicKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SSHPublicKeysClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	public_key := d.Get("public_key").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing SSH Public Key %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(resp.Response) {
		return tf.ImportAsExistsError("azurerm_ssh_public_key", *resp.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	t := d.Get("tags").(map[string]interface{})

	params := compute.SSHPublicKeyResource{
		Name:     utils.String(name),
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		SSHPublicKeyResourceProperties: &compute.SSHPublicKeyResourceProperties{
			PublicKey: utils.String(public_key),
		},
	}

	if _, err := client.Create(ctx, resourceGroup, name, params); err != nil {
		return fmt.Errorf("creating SSH Public Key %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving SSH Public Key %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("retrieving SSH Public Key %q (Resource Group %q): `id` was nil", name, resourceGroup)
	}

	d.SetId(*read.ID)
	return resourceSshPublicKeyRead(d, meta)
}

func resourceSshPublicKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SSHPublicKeysClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SSHPublicKeyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] SSH Public Key %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving SSH Public Key %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("id", resp.ID)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSshPublicKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SSHPublicKeysClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SSHPublicKeyID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Retrieving SSH Public Key %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving SSH Public Key %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	props := compute.SSHPublicKeyResourceProperties{}

	if d.HasChange("public_key") {
		props.PublicKey = utils.String(d.Get("public_key").(string))
	}

	update := compute.SSHPublicKeyUpdateResource{
		SSHPublicKeyResourceProperties: &props,
	}

	if d.HasChange("tags") {

		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	}

	log.Printf("[DEBUG] Updating SSH Public Key %q (Resource Group %q)..", id.Name, id.ResourceGroup)

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, update); err != nil {
		return fmt.Errorf("updating SSH Public Key %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Updated SSH Public Key %q (Resource Group %q).", id.Name, id.ResourceGroup)

	return resourceSshPublicKeyRead(d, meta)
}

func resourceSshPublicKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SSHPublicKeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SSHPublicKeyID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Retrieving SSH Public Key %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving SSH Public Key %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Deleting SSH Public Key %q (Resource Group %q)..", id.Name, id.ResourceGroup)
	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("deleting SSH Public Key %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Deleted SSH Public Key %q (Resource Group %q).", id.Name, id.ResourceGroup)

	return nil
}
