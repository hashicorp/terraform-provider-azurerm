package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/operationsmanagement/mgmt/2015-11-01-preview/operationsmanagement"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsSolution() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsSolutionCreateUpdate,
		Read:   resourceArmLogAnalyticsSolutionRead,
		Update: resourceArmLogAnalyticsSolutionCreateUpdate,
		Delete: resourceArmLogAnalyticsSolutionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameDiffSuppressSchema(),

			"plan": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"publisher": {
							Type:     schema.TypeString,
							Required: true,
						},
						"promotion_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"product": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"workspace_resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmLogAnalyticsSolutionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).solutionsClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics solution creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	workspaceID := d.Get("workspace_resource_id").(string)

	solutionPlan := expandAzureRmLogAnalyticsSolutionPlan(d)

	parameters := operationsmanagement.Solution{
		Name:     &name,
		Location: &location,
		Plan:     &solutionPlan,
		Properties: &operationsmanagement.SolutionProperties{
			WorkspaceResourceID: &workspaceID,
		},
	}

	solution, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	if solution.ID == nil {
		return fmt.Errorf("Cannot read Log Analytics Solution '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*solution.ID)

	return resourceArmLogAnalyticsSolutionRead(d, meta)

}

func resourceArmLogAnalyticsSolutionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).solutionsClient
	ctx := meta.(*ArmClient).StopContext
	id, err := parseAzureResourceID(d.Id())
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
		return fmt.Errorf("Error making Read request on AzureRM Log Analytics solutions '%s': %+v Plan was nil", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", resp.Location)
	d.Set("resource_group_name", resGroup)
	d.Set("workspace_resource_id", resp.ID)
	d.Set("plan", flattenAzureRmLogAnalyticsSolutionPlan(*resp.Plan))
	return nil
}

func resourceArmLogAnalyticsSolutionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).solutionsClient
	ctx := meta.(*ArmClient).StopContext
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["solutions"]

	resp, err := client.Delete(ctx, resGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Log Analytics Solution '%s': %+v", name, err)
	}

	return nil
}

func expandAzureRmLogAnalyticsSolutionPlan(d *schema.ResourceData) operationsmanagement.SolutionPlan {
	plans := d.Get("plan").([]interface{})
	plan := plans[0].(map[string]interface{})

	expandedPlan := operationsmanagement.SolutionPlan{}

	if name := plan["name"].(string); len(name) > 0 {
		expandedPlan.Name = &name
	}

	if publisher := plan["publisher"].(string); len(publisher) > 0 {
		expandedPlan.Publisher = &publisher
	}

	if promotionCode := plan["promotion_code"].(string); len(promotionCode) > 0 {
		expandedPlan.PromotionCode = &promotionCode
	} else {
		blankString := ""
		expandedPlan.PromotionCode = &blankString
	}

	if product := plan["product"].(string); len(product) > 0 {
		expandedPlan.Product = &product
	}

	return expandedPlan
}

func flattenAzureRmLogAnalyticsSolutionPlan(plan operationsmanagement.SolutionPlan) []interface{} {
	plans := make([]interface{}, 0)
	values := make(map[string]interface{})

	values["name"] = *plan.Name
	values["product"] = *plan.Product
	values["promotion_code"] = *plan.PromotionCode
	values["publisher"] = *plan.Publisher

	return append(plans, values)
}
