package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSecurityCenterAssessment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityCenterAssessmentCreateUpdate,
		Read:   resourceSecurityCenterAssessmentRead,
		Update: resourceSecurityCenterAssessmentCreateUpdate,
		Delete: resourceSecurityCenterAssessmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AssessmentID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"assessment_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AssessmentMetadataID,
			},

			"target_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"status": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(security.Healthy),
								string(security.NotApplicable),
								string(security.Unhealthy),
							}, false),
						},

						"cause": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"additional_data": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSecurityCenterAssessmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	metadataID, err := parse.AssessmentMetadataID(d.Get("assessment_policy_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewAssessmentID(d.Get("target_resource_id").(string), metadataID.AssessmentMetadataName)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.TargetResourceID, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Security Center Assessments %q : %+v", id.ID(), err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_security_center_assessment", id.ID())
		}
	}

	assessment := security.Assessment{
		AssessmentProperties: &security.AssessmentProperties{
			AdditionalData: utils.ExpandMapStringPtrString(d.Get("additional_data").(map[string]interface{})),
			ResourceDetails: &security.AzureResourceDetails{
				Source: security.SourceAzure,
			},
			Status: expandSecurityCenterAssessmentStatus(d.Get("status").([]interface{})),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.TargetResourceID, id.Name, assessment); err != nil {
		return fmt.Errorf("creating/updating Security Center Assessment %q (target resource id %q) : %+v", id.Name, id.TargetResourceID, err)
	}

	d.SetId(id.ID())

	return resourceSecurityCenterAssessmentRead(d, meta)
}

func resourceSecurityCenterAssessmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AssessmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.TargetResourceID, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] security Center Assessment %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Security Center Assessment %q (target resource id %q) : %+v", id.Name, id.TargetResourceID, err)
	}

	d.Set("assessment_policy_id", parse.NewAssessmentMetadataID(subscriptionID, id.Name).ID())
	d.Set("target_resource_id", id.TargetResourceID)
	if props := resp.AssessmentProperties; props != nil {
		d.Set("additional_data", utils.FlattenMapStringPtrString(props.AdditionalData))
		if err := d.Set("status", flattenSecurityCenterAssessmentStatus(props.Status)); err != nil {
			return fmt.Errorf("setting `status`: %s", err)
		}
	}

	return nil
}

func resourceSecurityCenterAssessmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AssessmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.TargetResourceID, id.Name); err != nil {
		return fmt.Errorf("deleting Security Center Assessment %q (target resource id %q) : %+v", id.Name, id.TargetResourceID, err)
	}

	return nil
}

func expandSecurityCenterAssessmentStatus(input []interface{}) *security.AssessmentStatus {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &security.AssessmentStatus{
		Code:        security.AssessmentStatusCode(v["code"].(string)),
		Cause:       utils.String(v["cause"].(string)),
		Description: utils.String(v["description"].(string)),
	}
}

func flattenSecurityCenterAssessmentStatus(input *security.AssessmentStatus) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var cause, description string
	if input.Cause != nil {
		cause = *input.Cause
	}
	if input.Description != nil {
		description = *input.Description
	}

	return []interface{}{
		map[string]interface{}{
			"code":        string(input.Code),
			"cause":       cause,
			"description": description,
		},
	}
}
