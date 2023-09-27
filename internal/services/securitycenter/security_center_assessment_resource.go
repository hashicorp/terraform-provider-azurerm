// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSecurityCenterAssessment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSecurityCenterAssessmentCreateUpdate,
		Read:   resourceSecurityCenterAssessmentRead,
		Update: resourceSecurityCenterAssessmentCreateUpdate,
		Delete: resourceSecurityCenterAssessmentDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := assessments.ParseScopedAssessmentID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"assessment_policy_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AssessmentMetadataID,
			},

			"target_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"status": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"code": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(assessments.AssessmentStatusCodeHealthy),
								string(assessments.AssessmentStatusCodeNotApplicable),
								string(assessments.AssessmentStatusCodeUnhealthy),
							}, false),
						},

						"cause": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"description": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"additional_data": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceSecurityCenterAssessmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	metadataID, err := assessmentsmetadata.ParseProviderAssessmentMetadataID(d.Get("assessment_policy_id").(string))
	if err != nil {
		return err
	}

	id := assessments.NewScopedAssessmentID(d.Get("target_resource_id").(string), metadataID.AssessmentMetadataName)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, assessments.GetOperationOptions{})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for present of existing Security Center Assessments %q : %+v", id.ID(), err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_security_center_assessment", id.ID())
		}
	}

	assessment := assessments.SecurityAssessment{
		Properties: &assessments.SecurityAssessmentProperties{
			AdditionalData: utils.ExpandPtrMapStringString(d.Get("additional_data").(map[string]interface{})),
			ResourceDetails: assessments.ResourceDetails{
				Source: assessments.SourceAzure,
			},
			Status: pointer.From(expandSecurityCenterAssessmentStatus(d.Get("status").([]interface{}))),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, assessment); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSecurityCenterAssessmentRead(d, meta)
}

func resourceSecurityCenterAssessmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assessments.ParseScopedAssessmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, assessments.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] security Center Assessment %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("assessment_policy_id", assessmentsmetadata.NewProviderAssessmentMetadataID(subscriptionID, id.AssessmentName).ID())
	d.Set("target_resource_id", id.ResourceId)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("additional_data", utils.FlattenPtrMapStringString(props.AdditionalData))
			if err := d.Set("status", flattenSecurityCenterAssessmentStatus(pointer.To(props.Status))); err != nil {
				return fmt.Errorf("setting `status`: %s", err)
			}
		}
	}

	return nil
}

func resourceSecurityCenterAssessmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assessments.ParseScopedAssessmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandSecurityCenterAssessmentStatus(input []interface{}) *assessments.AssessmentStatus {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &assessments.AssessmentStatus{
		Code:        assessments.AssessmentStatusCode(v["code"].(string)),
		Cause:       utils.String(v["cause"].(string)),
		Description: utils.String(v["description"].(string)),
	}
}

func flattenSecurityCenterAssessmentStatus(input *assessments.AssessmentStatusResponse) []interface{} {
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
