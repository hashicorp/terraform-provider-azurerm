package managementpartner

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementpartner/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmManagementPartner() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagementPartnerCreate,
		Read:   resourceArmManagementPartnerRead,
		Update: resourceArmManagementPartnerUpdate,
		Delete: resourceArmManagementPartnerDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ManagementPartnerID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"partner_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"partner_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmManagementPartnerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementPartner.PartnerClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	partnerId := d.Get("partner_id").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, partnerId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Management Partner %q: %+v", partnerId, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_management_partner", *existing.ID)
		}
	}

	_, err := client.Create(ctx, partnerId)
	if err != nil {
		return fmt.Errorf("Error creating Management Partner %q: %+v", partnerId, err)
	}

	resp, err := client.Get(ctx, partnerId)
	if err != nil {
		return fmt.Errorf("Error retrieving Management Partner %q: %+v", partnerId, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Management Partner %q ID", partnerId)
	}
	d.SetId(*resp.ID)

	return resourceArmManagementPartnerRead(d, meta)
}

func resourceArmManagementPartnerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementPartner.PartnerClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	partnerId := d.Get("partner_id").(string)

	_, err := client.Update(ctx, partnerId)
	if err != nil {
		return fmt.Errorf("Error updating Management Partner %q: %+v", partnerId, err)
	}

	return resourceArmManagementPartnerRead(d, meta)
}

func resourceArmManagementPartnerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementPartner.PartnerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagementPartnerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.PartnerId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Management Partner %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Management Partner %q: %+v", d.Id(), err)
	}

	d.Set("partner_id", resp.PartnerID)
	d.Set("partner_name", resp.PartnerName)

	return nil
}

func resourceArmManagementPartnerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementPartner.PartnerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagementPartnerID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.PartnerId)
	if err != nil {
		return fmt.Errorf("Error deleting Management Partner %q: %+v", id.PartnerId, err)
	}

	return nil
}
