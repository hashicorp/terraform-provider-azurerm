package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationsmanagement/mgmt/2015-11-01-preview/operationsmanagement"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	loganalyticsParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsSolution() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsSolutionCreateUpdate,
		Read:   resourceLogAnalyticsSolutionRead,
		Update: resourceLogAnalyticsSolutionCreateUpdate,
		Delete: resourceLogAnalyticsSolutionDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := loganalyticsParse.LogAnalyticsSolutionID(id)
			return err
		}),

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Log Analytics Solution creation.")

	// The resource requires both .name and .plan.name are set in the format
	// "SolutionName(WorkspaceName)". Feedback will be submitted to the OMS team as IMO this isn't ideal.
	id := loganalyticsParse.NewLogAnalyticsSolutionID(subscriptionId, d.Get("resource_group_name").(string), fmt.Sprintf("%s(%s)", d.Get("solution_name").(string), d.Get("workspace_name").(string)))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SolutionName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_log_analytics_solution", id.ID())
		}
	}

	solutionPlan := expandAzureRmLogAnalyticsSolutionPlan(d)
	solutionPlan.Name = &id.SolutionName

	location := azure.NormalizeLocation(d.Get("location").(string))
	workspaceID := d.Get("workspace_resource_id").(string)

	parameters := operationsmanagement.Solution{
		Name:     utils.String(id.SolutionName),
		Location: utils.String(location),
		Plan:     &solutionPlan,
		Properties: &operationsmanagement.SolutionProperties{
			WorkspaceResourceID: utils.String(workspaceID),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SolutionName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating Log Analytics Solution %q (Workspace %q / Resource Group %q): %+v", id.SolutionName, workspaceID, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the create/update of Log Analytics Solution %q (Workspace %q / Resource Group %q): %+v", id.SolutionName, workspaceID, id.ResourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceLogAnalyticsSolutionRead(d, meta)
}

func resourceLogAnalyticsSolutionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SolutionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := loganalyticsParse.LogAnalyticsSolutionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SolutionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	if resp.Plan == nil {
		return fmt.Errorf("making Read request on %s: Plan was nil", *id)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	// Reversing the mapping used to get .solution_name
	// expecting resp.Name to be in format "SolutionName(WorkspaceName)".
	if v := resp.Name; v != nil {
		val := *v
		segments := strings.Split(*v, "(")
		if len(segments) != 2 {
			return fmt.Errorf("expected %q to match 'Solution(WorkspaceName)'", val)
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
		return fmt.Errorf("setting `plan`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLogAnalyticsSolutionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SolutionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := loganalyticsParse.LogAnalyticsSolutionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SolutionName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
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
