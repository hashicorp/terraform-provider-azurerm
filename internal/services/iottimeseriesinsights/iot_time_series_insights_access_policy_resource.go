// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iottimeseriesinsights

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/accesspolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/environments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iottimeseriesinsights/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIoTTimeSeriesInsightsAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIoTTimeSeriesInsightsAccessPolicyCreateUpdate,
		Read:   resourceIoTTimeSeriesInsightsAccessPolicyRead,
		Update: resourceIoTTimeSeriesInsightsAccessPolicyCreateUpdate,
		Delete: resourceIoTTimeSeriesInsightsAccessPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := accesspolicies.ParseAccessPolicyID(id)
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
			0: migration.StandardEnvironmentAccessPolicyV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[-\w\._\(\)]+$`),
					"IoT Time Series Insights Access Policy name must contain only word characters, periods, underscores, hyphens, and parentheses.",
				),
			},

			"time_series_insights_environment_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: environments.ValidateEnvironmentID,
			},

			"principal_object_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"roles": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(accesspolicies.AccessPolicyRoleContributor),
						string(accesspolicies.AccessPolicyRoleReader),
					}, false),
				},
			},
		},
	}
}

func resourceIoTTimeSeriesInsightsAccessPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.AccessPolicies
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	environmentId, err := environments.ParseEnvironmentID(d.Get("time_series_insights_environment_id").(string))
	if err != nil {
		return err
	}

	id := accesspolicies.NewAccessPolicyID(subscriptionId, environmentId.ResourceGroupName, environmentId.EnvironmentName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_iot_time_series_insights_access_policy", id.ID())
		}
	}

	payload := accesspolicies.AccessPolicyCreateOrUpdateParameters{
		Properties: accesspolicies.AccessPolicyResourceProperties{
			Description:       utils.String(d.Get("description").(string)),
			PrincipalObjectId: utils.String(d.Get("principal_object_id").(string)),
			Roles:             expandIoTTimeSeriesInsightsAccessPolicyRoles(d.Get("roles").(*pluginsdk.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceIoTTimeSeriesInsightsAccessPolicyRead(d, meta)
}

func resourceIoTTimeSeriesInsightsAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.AccessPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accesspolicies.ParseAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AccessPolicyName)
	d.Set("time_series_insights_environment_id", environments.NewEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.EnvironmentName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("principal_object_id", props.PrincipalObjectId)
			if err := d.Set("roles", flattenIoTTimeSeriesInsightsAccessPolicyRoles(props.Roles)); err != nil {
				return fmt.Errorf("setting `roles`: %+v", err)
			}
		}
	}

	return nil
}

func resourceIoTTimeSeriesInsightsAccessPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.AccessPolicies
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := accesspolicies.ParseAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandIoTTimeSeriesInsightsAccessPolicyRoles(input []interface{}) *[]accesspolicies.AccessPolicyRole {
	roles := make([]accesspolicies.AccessPolicyRole, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		roles = append(roles, accesspolicies.AccessPolicyRole(v.(string)))
	}

	return &roles
}

func flattenIoTTimeSeriesInsightsAccessPolicyRoles(input *[]accesspolicies.AccessPolicyRole) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}
