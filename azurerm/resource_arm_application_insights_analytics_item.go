package azurerm

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArmApplicationInsightsAnalyticsItem() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationInsightsAnalyticsItemCreate,
		Read:   resourceArmApplicationInsightsAnalyticsItemRead,
		Update: resourceArmApplicationInsightsAnalyticsItemUpdate,
		Delete: resourceArmApplicationInsightsAnalyticsItemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"application_insights_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"item_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"content": {
				Type:     schema.TypeString,
				Required: true,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.ItemScopeShared),
					string(insights.ItemScopeUser),
				}, false),
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.Query),
					string(insights.Function),
					string(insights.Folder),
					string(insights.Recent),
				}, false),
			},

			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"time_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmApplicationInsightsAnalyticsItemCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmApplicationInsightsAnalyticsItemCreateUpdate(d, meta, false)
}
func resourceArmApplicationInsightsAnalyticsItemUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmApplicationInsightsAnalyticsItemCreateUpdate(d, meta, true)
}
func resourceArmApplicationInsightsAnalyticsItemCreateUpdate(d *schema.ResourceData, meta interface{}, overwrite bool) error {
	return fmt.Errorf("Not implemented")
}

func resourceArmApplicationInsightsAnalyticsItemRead(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Not implemented")
}

func resourceArmApplicationInsightsAnalyticsItemDelete(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("Not implemented")
}
