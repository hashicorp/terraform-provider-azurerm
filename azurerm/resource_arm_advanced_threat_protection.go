package azurerm

import (
	"fmt"
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
		Delete:   resourceArmAdvancedThreatProtectionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enable": {
				Type:         schema.TypeBool,
				Required:     true,
			},
		},
	}
}

func resourceArmAdvancedThreatProtectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceID := d.Get("resource_id").(string)

	setting := security.AdvancedThreatProtectionSetting{
		AdvancedThreatProtectionProperties: &security.AdvancedThreatProtectionProperties{
			IsEnabled: utils.Bool(d.Get("enable").(bool)),
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

	id := d.Id()
	parts := strings.Split(strings.Trim(id, "/"), "/providers/Microsoft.Security/advancedThreatProtectionSettings/")
	if len(parts) != 2 {
		return fmt.Errorf("Error determining target resource ID, resource ID in unexacpted format: %q", id)
	}
	resourceID := parts[0]

	atp, err := client.Get(ctx, resourceID)
	if err != nil {
		return fmt.Errorf("Error reading Advanced Threat protection for resource %q: %+v", resourceID, err)
	}

	if atpp := atp.AdvancedThreatProtectionProperties; atpp != nil {
		d.Set("enable", atpp.IsEnabled)
	}

	return nil
}

func resourceArmAdvancedThreatProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id := d.Id()
	parts := strings.Split(strings.Trim(id, "/"), "/providers/Microsoft.Security/advancedThreatProtectionSettings/")
	if len(parts) != 2 {
		return fmt.Errorf("Error determining target resource ID, resource ID in unexacpted format: %q", id)
	}
	resourceID := parts[0]

	setting := security.AdvancedThreatProtectionSetting{
		AdvancedThreatProtectionProperties: &security.AdvancedThreatProtectionProperties{
			IsEnabled: utils.Bool(false),
		},
	}

	if _, err := client.Create(ctx, resourceID, setting); err != nil {
		return fmt.Errorf("Error resetting Advanced Threat protection for resource %q: %+v", resourceID, err)
	}

	return nil
}