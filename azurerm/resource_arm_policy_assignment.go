package azurerm

import (
	"fmt"
	"log"

	"time"

	"context"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-12-01/policy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicyAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPolicyAssignmentCreate,
		Read:   resourceArmPolicyAssignmentRead,
		Delete: resourceArmPolicyAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"policy_definition_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"parameters": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
		},
	}
}

func resourceArmPolicyAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).policyAssignmentsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	policyDefinitionId := d.Get("policy_definition_id").(string)
	displayName := d.Get("display_name").(string)

	assignment := policy.Assignment{
		AssignmentProperties: &policy.AssignmentProperties{
			PolicyDefinitionID: utils.String(policyDefinitionId),
			DisplayName:        utils.String(displayName),
			Scope:              utils.String(scope),
		},
	}

	if v := d.Get("description").(string); v != "" {
		assignment.AssignmentProperties.Description = utils.String(v)
	}

	if v := d.Get("parameters").(string); v != "" {
		expandedParams, err := structure.ExpandJsonFromString(v)
		if err != nil {
			return fmt.Errorf("Error expanding JSON from Parameters %q: %+v", v, err)
		}

		assignment.AssignmentProperties.Parameters = &expandedParams
	}

	_, err := client.Create(ctx, scope, name, assignment)
	if err != nil {
		return err
	}

	// Policy Assignments are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Assignment %q to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policyAssignmentRefreshFunc(ctx, client, scope, name),
		Timeout:                   5 * time.Minute,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Policy Assignment %q to become available: %s", name, err)
	}

	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmPolicyAssignmentRead(d, meta)
}

func resourceArmPolicyAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).policyAssignmentsClient
	ctx := meta.(*ArmClient).StopContext

	id := d.Id()

	resp, err := client.GetByID(ctx, id)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Policy Assignment %q - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Policy Assignment %q: %+v", id, err)
	}

	d.Set("name", resp.Name)

	if props := resp.AssignmentProperties; props != nil {
		d.Set("scope", props.Scope)
		d.Set("policy_definition_id", props.PolicyDefinitionID)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)

		if params := props.Parameters; params != nil {
			paramsVal := params.(map[string]interface{})
			json, err := structure.FlattenJsonToString(paramsVal)
			if err != nil {
				return fmt.Errorf("Error serializing JSON from Parameters: %+v", err)
			}

			d.Set("parameters", json)
		}
	}

	return nil
}

func resourceArmPolicyAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).policyAssignmentsClient
	ctx := meta.(*ArmClient).StopContext

	id := d.Id()

	resp, err := client.DeleteByID(ctx, id)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("Error deleting Policy Assignment %q: %+v", id, err)
	}

	return nil
}

func policyAssignmentRefreshFunc(ctx context.Context, client policy.AssignmentsClient, scope string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, scope, name)
		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("Error issuing read request in policyAssignmentRefreshFunc for Policy Assignment %q (Scope: %q): %s", name, scope, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
