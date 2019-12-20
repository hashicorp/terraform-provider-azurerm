package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter"
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

func resourceArmAdvancedThreatProtectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceID := d.Get("target_resource_id").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		server, err := client.Get(ctx, resourceID)
		if err != nil {
			if !utils.ResponseWasNotFound(server.Response) {
				return fmt.Errorf("Error checking for presence of existing Advanced Threat Protection for resource %q: %+v", resourceID, err)
			}
		}

		if server.ID != nil && *server.ID != "" && server.IsEnabled != nil && *server.IsEnabled {
			return tf.ImportAsExistsError("azurerm_advanced_threat_protection", *server.ID)
		}
	}

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
	client := meta.(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitycenter.ParseAdvancedThreatProtectionID(d.Id())
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
	client := meta.(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitycenter.ParseAdvancedThreatProtectionID(d.Id())
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
