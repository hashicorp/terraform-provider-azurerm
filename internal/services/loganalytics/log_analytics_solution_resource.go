// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationsmanagement/2015-11-01-preview/solution"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
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
			_, err := solution.ParseSolutionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SolutionV0ToV1{},
		}),

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

			"location": commonschema.Location(),

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
	id := solution.NewSolutionID(subscriptionId, d.Get("resource_group_name").(string), fmt.Sprintf("%s(%s)", d.Get("solution_name").(string), d.Get("workspace_name").(string)))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_log_analytics_solution", id.ID())
		}
	}

	solutionPlan := expandAzureRmLogAnalyticsSolutionPlan(d)
	solutionPlan.Name = &id.SolutionName

	location := azure.NormalizeLocation(d.Get("location").(string))
	workspaceID, err := workspaces.ParseWorkspaceID(d.Get("workspace_resource_id").(string))
	if err != nil {
		return err
	}

	parameters := solution.Solution{
		Name:     utils.String(id.SolutionName),
		Location: utils.String(location),
		Plan:     &solutionPlan,
		Properties: &solution.SolutionProperties{
			WorkspaceResourceId: workspaceID.ID(),
		},
		Tags: expandTags(d.Get("tags").(map[string]interface{})),
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLogAnalyticsSolutionRead(d, meta)
}

func resourceLogAnalyticsSolutionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SolutionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := solution.ParseSolutionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if model.Plan == nil {
			return fmt.Errorf("making Read request on %s: Plan was nil", *id)
		}

		d.Set("resource_group_name", id.ResourceGroupName)
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		// Reversing the mapping used to get .solution_name
		// expecting resp.Name to be in format "SolutionName(WorkspaceName)".
		if v := model.Name; v != nil {
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

		if props := model.Properties; props != nil {
			var workspaceId string
			if props.WorkspaceResourceId != "" {
				id, err := workspaces.ParseWorkspaceIDInsensitively(props.WorkspaceResourceId)
				if err != nil {
					return err
				}
				workspaceId = id.ID()
			}
			d.Set("workspace_resource_id", workspaceId)
		}

		if err := d.Set("plan", flattenAzureRmLogAnalyticsSolutionPlan(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}

		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceLogAnalyticsSolutionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SolutionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := solution.ParseSolutionID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandAzureRmLogAnalyticsSolutionPlan(d *pluginsdk.ResourceData) solution.SolutionPlan {
	plans := d.Get("plan").([]interface{})
	plan := plans[0].(map[string]interface{})

	name := plan["name"].(string)
	publisher := plan["publisher"].(string)
	promotionCode := plan["promotion_code"].(string)
	product := plan["product"].(string)

	expandedPlan := solution.SolutionPlan{
		Name:          utils.String(name),
		PromotionCode: utils.String(promotionCode),
		Publisher:     utils.String(publisher),
		Product:       utils.String(product),
	}

	return expandedPlan
}

func flattenAzureRmLogAnalyticsSolutionPlan(input *solution.SolutionPlan) []interface{} {
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
