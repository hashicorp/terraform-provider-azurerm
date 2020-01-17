package managementpartner

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmManagementPartner() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmManagementPartnerRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
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

func dataSourceArmManagementPartnerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagementPartner.PartnerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	partnerId := d.Get("partner_id").(string)

	resp, err := client.Get(ctx, partnerId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Management Partner %q was not found", partnerId)
		}
		return fmt.Errorf("Error reading Management Partner %q : %+v", partnerId, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Error retrieving Management Partner %q : ID was nil or empty", partnerId)
	}

	d.SetId(*resp.ID)
	d.Set("partner_id", partnerId)
	d.Set("partner_name", (resp.PartnerName))

	return nil
}
