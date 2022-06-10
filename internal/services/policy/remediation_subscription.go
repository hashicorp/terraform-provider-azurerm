package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	validate2 "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/sdk/2021-10-01/policyinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmSubscriptionPolicyRemediation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmSubscriptionPolicyRemediationCreateUpdate,
		Read:   resourceArmSubscriptionPolicyRemediationRead,
		Update: resourceArmSubscriptionPolicyRemediationCreateUpdate,
		Delete: resourceArmSubscriptionPolicyRemediationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SubscriptionPolicyRemediationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RemediationName,
			},

			"subscription_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubscriptionID,
			},

			"policy_assignment_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				// TODO: remove this suppression when github issue https://github.com/Azure/azure-rest-api-specs/issues/8353 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.PolicyAssignmentID,
			},

			"failure_percentage": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Description:  "A number between 0.0 to 1.0 representing the percentage failure threshold.",
				ValidateFunc: validate2.FloatInRange(0, 1.0),
			},

			"parallel_deployments": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validate2.IntegerPositive,
			},

			"resource_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validate2.IntegerPositive,
			},

			"location_filters": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: location.EnhancedValidate,
				},
			},

			"policy_definition_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// TODO: remove this suppression when github issue https://github.com/Azure/azure-rest-api-specs/issues/8353 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.PolicyDefinitionID,
			},

			"resource_discovery_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(policyinsights.ResourceDiscoveryModeExistingNonCompliant),
				ValidateFunc: validation.StringInSlice([]string{
					string(policyinsights.ResourceDiscoveryModeExistingNonCompliant),
					string(policyinsights.ResourceDiscoveryModeReEvaluateCompliance),
				}, false),
			},
		},
	}
}

func resourceArmSubscriptionPolicyRemediationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicyInsightsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId, err := commonids.ParseSubscriptionID(d.Get("subscription_id").(string))
	if err != nil {
		return err
	}

	id := policyinsights.NewRemediationID(subscriptionId.SubscriptionId, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.RemediationsGetAtSubscription(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
			}
		}
		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_subscription_policy_remediation", *existing.Model.Id)
		}
	}

	parameters := policyinsights.Remediation{
		Properties: &policyinsights.RemediationProperties{
			Filters: &policyinsights.RemediationFilters{
				Locations: utils.ExpandStringSlice(d.Get("location_filters").([]interface{})),
			},
			PolicyAssignmentId:          utils.String(d.Get("policy_assignment_id").(string)),
			PolicyDefinitionReferenceId: utils.String(d.Get("policy_definition_id").(string)),
		},
	}
	if v := d.Get("resource_discovery_mode").(string); v != "" {
		mode := policyinsights.ResourceDiscoveryMode(v)
		parameters.Properties.ResourceDiscoveryMode = &mode
	}
	if v := d.Get("failure_percentage").(float64); v != 0 {
		parameters.Properties.FailureThreshold = &policyinsights.RemediationPropertiesFailureThreshold{
			Percentage: utils.Float(v),
		}
	}
	if v := d.Get("parallel_deployments").(int); v != 0 {
		parameters.Properties.ParallelDeployments = utils.Int64(int64(v))
	}
	if v := d.Get("resource_count").(int); v != 0 {
		parameters.Properties.ResourceCount = utils.Int64(int64(v))
	}

	if _, err = client.RemediationsCreateOrUpdateAtSubscription(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id.ID(), err)
	}

	d.SetId(id.ID())

	return resourceArmSubscriptionPolicyRemediationRead(d, meta)
}

func resourceArmSubscriptionPolicyRemediationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicyInsightsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyinsights.ParseRemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("reading Policy Remediation: %+v", err)
	}

	subscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	resp, err := client.RemediationsGetAtSubscription(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id.ID(), err)
	}

	d.Set("name", id.RemediationName)
	d.Set("subscription_id", subscriptionId.ID())

	if props := resp.Model.Properties; props != nil {
		locations := []interface{}{}
		if filters := props.Filters; filters != nil {
			locations = utils.FlattenStringSlice(filters.Locations)
		}
		if err := d.Set("location_filters", locations); err != nil {
			return fmt.Errorf("setting `location_filters`: %+v", err)
		}

		d.Set("policy_assignment_id", props.PolicyAssignmentId)
		d.Set("policy_definition_id", props.PolicyDefinitionReferenceId)
		d.Set("resource_discovery_mode", utils.NormalizeNilableString((*string)(props.ResourceDiscoveryMode)))

		d.Set("resource_count", props.ResourceCount)
		d.Set("parallel_deployments", props.ParallelDeployments)
		if props.FailureThreshold != nil {
			d.Set("failure_percentage", props.FailureThreshold.Percentage)
		}
	}

	return nil
}

func resourceArmSubscriptionPolicyRemediationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicyInsightsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyinsights.ParseRemediationID(d.Id())
	if err != nil {
		return err
	}

	// we have to cancel the remediation first before deleting it when the resource_discovery_mode is set to ReEvaluateCompliance
	// therefore we first retrieve the remediation to see if the resource_discovery_mode is switched to ReEvaluateCompliance
	existing, err := client.RemediationsGetAtSubscription(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if prop := existing.Model.Properties; prop != nil {
		if mode := existing.Model.Properties.ResourceDiscoveryMode; mode != nil && *mode == policyinsights.ResourceDiscoveryModeReEvaluateCompliance {
			// Remediation can only be canceld when it is in "Evaluating" status, otherwise, API might raise error (e.g. canceling a "Completed" remediation returns 400).
			if existing.Model.Properties.ProvisioningState != nil && *existing.Model.Properties.ProvisioningState == "Evaluating" {
				log.Printf("[DEBUG] cancelling the remediation first before deleting it when `resource_discovery_mode` is set to `ReEvaluateCompliance`")
				if _, err := client.RemediationsCancelAtSubscription(ctx, *id); err != nil {
					return fmt.Errorf("cancelling %s: %+v", id.ID(), err)
				}

				log.Printf("[DEBUG] waiting for the %s to be canceled", id.ID())
				stateConf := &pluginsdk.StateChangeConf{
					Pending: []string{"Cancelling"},
					Target: []string{
						"Succeeded", "Canceled", "Failed",
					},
					Refresh:    subscriptionPolicyRemediationCancellationRefreshFunc(ctx, client, *id),
					MinTimeout: 10 * time.Second,
					Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
				}

				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("waiting for %s to be canceled: %+v", id.ID(), err)
				}
			}
		}
	}

	_, err = client.RemediationsDeleteAtSubscription(ctx, *id)

	return err
}

func subscriptionPolicyRemediationCancellationRefreshFunc(ctx context.Context, client *policyinsights.PolicyInsightsClient, id policyinsights.RemediationId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.RemediationsGetAtSubscription(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("issuing read request for %s: %+v", id.ID(), err)
		}

		if resp.Model.Properties == nil {
			return nil, "", fmt.Errorf("`properties` was nil")
		}
		if resp.Model.Properties.ProvisioningState == nil {
			return nil, "", fmt.Errorf("`properties.ProvisioningState` was nil")
		}
		return resp, *resp.Model.Properties.ProvisioningState, nil
	}
}
