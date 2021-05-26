package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationsmanagement/mgmt/2015-11-01-preview/operationsmanagement"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceLogAnalyticsSolution() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsSolutionCreateUpdate,
		Read:   resourceLogAnalyticsSolutionRead,
		Update: resourceLogAnalyticsSolutionCreateUpdate,
		Delete: resourceLogAnalyticsSolutionDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"solution_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"workspace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceName,
			},

			"workspace_resource_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"plan": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"publisher": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
						},
						"promotion_code": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"product": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceLogAnalyticsSolutionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SolutionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Log Analytics Solution creation.")

	// The resource requires both .name and .plan.name are set in the format
	// "SolutionName(WorkspaceName)". Feedback will be submitted to the OMS team as IMO this isn't ideal.
	name := fmt.Sprintf("%s(%s)", d.Get("solution_name").(string), d.Get("workspace_name").(string))
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Log Analytics Solution %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_solution", *existing.ID)
		}
	}

	solutionPlan := expandAzureRmLogAnalyticsSolutionPlan(d)
	solutionPlan.Name = &name

	location := azure.NormalizeLocation(d.Get("location").(string))
	workspaceID := d.Get("workspace_resource_id").(string)

	parameters := operationsmanagement.Solution{
		Name:     utils.String(name),
		Location: utils.String(location),
		Plan:     &solutionPlan,
		Properties: &operationsmanagement.SolutionProperties{
			WorkspaceResourceID: utils.String(workspaceID),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Log Analytics Solution %q (Workspace %q / Resource Group %q): %+v", name, workspaceID, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the create/update of Log Analytics Solution %q (Workspace %q / Resource Group %q): %+v", name, workspaceID, resGroup, err)
	}

	solution, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Log Analytics Solution %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if solution.ID == nil {
		return fmt.Errorf("Cannot read Log Analytics Solution %q (Resource Group %q) ID", name, resGroup)
	}

	d.SetId(*solution.ID)

	return resourceLogAnalyticsSolutionRead(d, meta)
}

func resourceLogAnalyticsSolutionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SolutionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["solutions"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Log Analytics solutions '%s': %+v", name, err)
	}

	if resp.Plan == nil {
		return fmt.Errorf("Error making Read request on AzureRM Log Analytics solutions '%s': Plan was nil", name)
	}

	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	// Reversing the mapping used to get .solution_name
	// expecting resp.Name to be in format "SolutionName(WorkspaceName)".
	if v := resp.Name; v != nil {
		val := *v
		segments := strings.Split(*v, "(")
		if len(segments) != 2 {
			return fmt.Errorf("Expected %q to match 'Solution(WorkspaceName)'", val)
		}

		solutionName := segments[0]
		workspaceName := strings.TrimSuffix(segments[1], ")")
		d.Set("solution_name", solutionName)
		d.Set("workspace_name", workspaceName)
	}

	if props := resp.Properties; props != nil {
		var workspaceId string
		if props.WorkspaceResourceID != nil {
			id, err := loganalyticsParse.LogAnalyticsWorkspaceID(*props.WorkspaceResourceID)
			if err != nil {
				return err
			}
			workspaceId = id.ID()
		}
		d.Set("workspace_resource_id", workspaceId)
	}

	if err := d.Set("plan", flattenAzureRmLogAnalyticsSolutionPlan(resp.Plan)); err != nil {
		return fmt.Errorf("Error setting `plan`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLogAnalyticsSolutionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SolutionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["solutions"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Log Analytics Solution %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Log Analytics Solution %q (Resource Group %q): %+v", name, resGroup, err)
		}
	}

	return nil
}

func expandAzureRmLogAnalyticsSolutionPlan(d *pluginsdk.ResourceData) operationsmanagement.SolutionPlan {
	plans := d.Get("plan").([]interface{})
	plan := plans[0].(map[string]interface{})

	name := plan["name"].(string)
	publisher := plan["publisher"].(string)
	promotionCode := plan["promotion_code"].(string)
	product := plan["product"].(string)

	expandedPlan := operationsmanagement.SolutionPlan{
		Name:          utils.String(name),
		PromotionCode: utils.String(promotionCode),
		Publisher:     utils.String(publisher),
		Product:       utils.String(product),
	}

	return expandedPlan
}

func flattenAzureRmLogAnalyticsSolutionPlan(input *operationsmanagement.SolutionPlan) []interface{} {
	output := make([]interface{}, 0)
	if input == nil {
		return output
	}

	values := make(map[string]interface{})

	if input.Name != nil {
		values["name"] = *input.Name
	}

	if input.Product != nil {
		values["product"] = *input.Product
	}

	if input.PromotionCode != nil {
		values["promotion_code"] = *input.PromotionCode
	}

	if input.Publisher != nil {
		values["publisher"] = *input.Publisher
	}

	return append(output, values)
}
