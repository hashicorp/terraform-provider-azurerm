// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package advisor

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2020-01-01/getrecommendations" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceAdvisorRecommendations() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAdvisorRecommendationsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"filter_by_category": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(getrecommendations.CategoryHighAvailability),
						string(getrecommendations.CategorySecurity),
						string(getrecommendations.CategoryPerformance),
						string(getrecommendations.CategoryCost),
						string(getrecommendations.CategoryOperationalExcellence),
					}, false),
				},
			},

			"filter_by_resource_groups": commonschema.ResourceGroupNameSetOptional(),

			"recommendations": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"category": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"impact": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"recommendation_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"recommendation_type_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"resource_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"resource_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"suppression_names": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"updated_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAdvisorRecommendationsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Advisor.RecommendationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSubscriptionID(meta.(*clients.Client).Account.SubscriptionId)

	filterList := make([]string, 0)
	if categories := expandAzureRmAdvisorRecommendationsMapString("Category", d.Get("filter_by_category").(*pluginsdk.Set).List()); categories != "" {
		filterList = append(filterList, categories)
	}
	if resGroups := expandAzureRmAdvisorRecommendationsMapString("ResourceGroup", d.Get("filter_by_resource_groups").(*pluginsdk.Set).List()); resGroups != "" {
		filterList = append(filterList, resGroups)
	}

	opts := getrecommendations.RecommendationsListOperationOptions{}
	if len(filterList) > 0 {
		opts.Filter = pointer.To(strings.Join(filterList, " and "))
	}

	recomendations, err := client.RecommendationsListComplete(ctx, id, opts)
	if err != nil {
		return fmt.Errorf("loading Advisor Recommendation for %q: %+v", id, err)
	}

	if err := d.Set("recommendations", flattenAzureRmAdvisorRecommendations(recomendations.Items)); err != nil {
		return fmt.Errorf("setting `recommendations`: %+v", err)
	}

	d.SetId(fmt.Sprintf("avdisor/recommendations/%s", time.Now().UTC().String()))

	return nil
}

func flattenAzureRmAdvisorRecommendations(recommends []getrecommendations.ResourceRecommendationBase) []interface{} {
	result := make([]interface{}, 0)

	if len(recommends) == 0 {
		return result
	}

	for _, r := range recommends {
		var description string
		var suppressionIds []interface{}

		v := r.Properties

		if v.ShortDescription != nil && v.ShortDescription.Problem != nil {
			description = *v.ShortDescription.Problem
		}

		if v.SuppressionIds != nil {
			suppressionIds = flattenSuppressionSlice(v.SuppressionIds)
		}

		result = append(result, map[string]interface{}{
			"category":               string(pointer.From(v.Category)),
			"description":            description,
			"impact":                 string(pointer.From(v.Impact)),
			"recommendation_name":    pointer.From(r.Name),
			"recommendation_type_id": pointer.From(v.RecommendationTypeId),
			"resource_name":          pointer.From(v.ImpactedValue),
			"resource_type":          pointer.From(v.ImpactedField),
			"suppression_names":      suppressionIds,
			"updated_time":           pointer.From(v.LastUpdated),
		})
	}

	return result
}

func expandAzureRmAdvisorRecommendationsMapString(t string, input []interface{}) string {
	if len(input) == 0 {
		return ""
	}
	result := make([]string, 0)
	for _, v := range input {
		result = append(result, fmt.Sprintf("%s eq '%s'", t, v.(string)))
	}
	return "(" + strings.Join(result, " or ") + ")"
}

func flattenSuppressionSlice(input *[]string) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}
