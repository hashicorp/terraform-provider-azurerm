package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAdvancedThreatProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAdvancedThreatProtectionCreateUpdate,
		Read:   resourceArmAdvancedThreatProtectionRead,
		Update: resourceArmAdvancedThreatProtectionCreateUpdate,
		Delete: resourceArmAdvancedThreatProtectionDelete,

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
			"target_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

type AdvancedThreatProtectionResourceID struct {
	Base azure.ResourceID

	TargetResourceID string
}

func parseAdvancedThreatProtectionID(input string) (*AdvancedThreatProtectionResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Advanced Threat Protection Set ID %q: %+v", input, err)
	}

	parts := strings.Split(input, "/providers/Microsoft.Security/advancedThreatProtectionSettings/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Error determining target resource ID, resource ID in unexacpted format: %q", id)
	}

	return &AdvancedThreatProtectionResourceID{
		Base:             *id,
		TargetResourceID: parts[0],
	}, nil

}

func resourceArmAdvancedThreatProtectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceID := d.Get("target_resource_id").(string)

	setting := security.AdvancedThreatProtectionSetting{
		AdvancedThreatProtectionProperties: &security.AdvancedThreatProtectionProperties{
			IsEnabled: utils.Bool(d.Get("enabled").(bool)),
		},
	}

	resp, err := client.Create(ctx, resourceID, setting)
	if err != nil {
		return fmt.Errorf("Error updating Advanced Threat protection for resource %q: %+v", resourceID, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Advanced Threat Protection for resource %q ", resourceID)
	}
	d.SetId(*resp.ID)

	return resourceArmAdvancedThreatProtectionRead(d, meta)
}

func resourceArmAdvancedThreatProtectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := parseAdvancedThreatProtectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.TargetResourceID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("Advanced Threat Protection was not found for resource %q: %+v", id.TargetResourceID, err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Advanced Threat protection for resource %q: %+v", id.TargetResourceID, err)
	}

	d.Set("target_resource_id", id.TargetResourceID)
	if atpp := resp.AdvancedThreatProtectionProperties; atpp != nil {
		d.Set("enabled", resp.IsEnabled)
	}

	return nil
}

func resourceArmAdvancedThreatProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := parseAdvancedThreatProtectionID(d.Id())
	if err != nil {
		return err
	}

	// there is no delete.. so lets just do best effort and set it to false?
	setting := security.AdvancedThreatProtectionSetting{
		AdvancedThreatProtectionProperties: &security.AdvancedThreatProtectionProperties{
			IsEnabled: utils.Bool(false),
		},
	}

	if _, err := client.Create(ctx, id.TargetResourceID, setting); err != nil {
		return fmt.Errorf("Error resetting Advanced Threat protection to false for resource %q: %+v", id.TargetResourceID, err)
	}

	return nil
}
