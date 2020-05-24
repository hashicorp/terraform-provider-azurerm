package loganalytics

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsWorkspaceCreateUpdate,
		Read:   resourceArmLogAnalyticsWorkspaceRead,
		Update: resourceArmLogAnalyticsWorkspaceCreateUpdate,
		Delete: resourceArmLogAnalyticsWorkspaceDelete,
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateAzureRmLogAnalyticsWorkspaceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(operationalinsights.Free),
					string(operationalinsights.PerGB2018),
					string(operationalinsights.PerNode),
					string(operationalinsights.Premium),
					string(operationalinsights.Standalone),
					string(operationalinsights.Standard),
					"Unlimited", // TODO check if this is actually no longer valid, removed in v28.0.0 of the SDK
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"retention_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.Any(validation.IntBetween(30, 730), validation.IntInSlice([]int{7})),
			},

			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"portal_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_shared_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_shared_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmLogAnalyticsWorkspaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Log Analytics Workspace %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_workspace", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	skuName := d.Get("sku").(string)
	sku := &operationalinsights.Sku{
		Name: operationalinsights.SkuNameEnum(skuName),
	}

	retentionInDays := int32(d.Get("retention_in_days").(int))

	t := d.Get("tags").(map[string]interface{})

	parameters := operationalinsights.Workspace{
		Name:     &name,
		Location: &location,
		Tags:     tags.Expand(t),
		WorkspaceProperties: &operationalinsights.WorkspaceProperties{
			Sku:             sku,
			RetentionInDays: &retentionInDays,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Log Analytics Workspace '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmLogAnalyticsWorkspaceRead(d, meta)
}

func resourceArmLogAnalyticsWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.LogAnalyticsWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Log Analytics workspaces '%s': %+v", id.Name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("workspace_id", resp.CustomerID)
	d.Set("portal_url", resp.PortalURL)
	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}
	d.Set("retention_in_days", resp.RetentionInDays)

	sharedKeys, err := client.GetSharedKeys(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		log.Printf("[ERROR] Unable to List Shared keys for Log Analytics workspaces %s: %+v", id.Name, err)
	} else {
		d.Set("primary_shared_key", sharedKeys.PrimarySharedKey)
		d.Set("secondary_shared_key", sharedKeys.SecondarySharedKey)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmLogAnalyticsWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.LogAnalyticsWorkspaceID(d.Id())
	if err != nil {
		return err
	}
	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Log Analytics Workspaces '%s': %+v", id.Name, err)
	}

	return nil
}

func ValidateAzureRmLogAnalyticsWorkspaceName(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile("^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$").MatchString(value) {
		errors = append(errors, fmt.Errorf("Workspace Name can only contain alphabet, number, and '-' character. You can not use '-' as the start and end of the name"))
	}

	length := len(value)
	if length > 63 || 4 > length {
		errors = append(errors, fmt.Errorf("Workspace Name can only be between 4 and 63 letters"))
	}

	return warnings, errors
}
