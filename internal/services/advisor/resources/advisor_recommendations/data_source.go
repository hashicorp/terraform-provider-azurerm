// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package advisor_recommendations

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/getrecommendations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Advisor Recommendations.",
		Read:        dataSourceAdvisorRecommendationsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"filter_by_category": {
				Description: "Specifies a list of categories in which the Advisor Recommendations will be listed.",
				Type:        pluginsdk.TypeSet,
				Optional:    true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice(getrecommendations.PossibleValuesForCategory(), false),
				},
			},

			"filter_by_resource_groups": commonschema.ResourceGroupNameSetOptional(),

			"filter_by_resource_ids": {
				Description: "Specifies a list of resource about which the Advisor Recommendations will be listed.",
				Type:        pluginsdk.TypeSet,
				Optional:    true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"filter_by_recommendation_type_guids": {
				Description: "Specifies a list of recommendation types about which the Advisor Recommendations will be listed.",
				Type:        pluginsdk.TypeSet,
				Optional:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"recommendations": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"category": {
							Description: "The category of the recommendation.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"description": {
							Description: "The description of the issue or the opportunity identified by the recommendation.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"id": {
							Description: "The name of the Advisor Recommendation.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"impact": {
							Description: "The business impact of the recommendation.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"recommendation_name": {
							Description: "The name of the Advisor Recommendation.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"recommendation_type_id": {
							Description: "The recommendation type id of the Advisor Recommendation.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"resource_name": {
							Description: "The name of the identified resource of the Advisor Recommendation.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"resource_type": {
							Description: "The type of the identified resource of the Advisor Recommendation.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"suppression_names": {
							Description: "A list of Advisor Suppression names of the Advisor Recommendation.",
							Type:        pluginsdk.TypeSet,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"updated_time": {
							Description: "The most recent time that Advisor checked the validity of the recommendation..",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAdvisorRecommendationsRead(d *pluginsdk.ResourceData, meta any) error {
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
	if resIDs := expandAzureRmAdvisorRecommendationsMapString("ResourceId", d.Get("filter_by_resource_ids").(*pluginsdk.Set).List()); resIDs != "" {
		filterList = append(filterList, resIDs)
	}
	if recommendationTypeGuids := expandAzureRmAdvisorRecommendationsMapString("RecommendationTypeGuid", d.Get("filter_by_recommendation_type_guids").(*pluginsdk.Set).List()); recommendationTypeGuids != "" {
		filterList = append(filterList, recommendationTypeGuids)
	}

	opts := getrecommendations.RecommendationsListOperationOptions{}
	if len(filterList) > 0 {
		opts.Filter = new(strings.Join(filterList, " and "))
	}

	recommendations, err := client.RecommendationsListComplete(ctx, id, opts)
	if err != nil {
		return fmt.Errorf("loading Advisor Recommendation for %q: %+v", id, err)
	}

	if err := d.Set("recommendations", flattenAzureRmAdvisorRecommendations(recommendations.Items)); err != nil {
		return fmt.Errorf("setting `recommendations`: %+v", err)
	}

	d.SetId(fmt.Sprintf("advisor/recommendations/%s/%d", id.SubscriptionId, pluginsdk.HashString(strings.Join(filterList, " and "))))

	return nil
}

func flattenAzureRmAdvisorRecommendations(recommends []getrecommendations.ResourceRecommendationBase) []any {
	result := make([]any, 0)

	if len(recommends) == 0 {
		return result
	}

	for _, r := range recommends {
		var description string
		var suppressionIds []any

		v := r.Properties

		if v.ShortDescription != nil && v.ShortDescription.Problem != nil {
			description = *v.ShortDescription.Problem
		}

		if v.SuppressionIds != nil {
			suppressionIds = flattenSuppressionSlice(v.SuppressionIds)
		}

		result = append(result, map[string]any{
			"category":               pointer.FromEnum(v.Category),
			"description":            description,
			"id":                     pointer.From(r.Id),
			"impact":                 pointer.FromEnum(v.Impact),
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

func expandAzureRmAdvisorRecommendationsMapString(t string, input []any) string {
	if len(input) == 0 {
		return ""
	}
	result := make([]string, 0)
	for _, v := range input {
		result = append(result, fmt.Sprintf("%s eq '%s'", t, v.(string)))
	}
	return "(" + strings.Join(result, " or ") + ")"
}

func flattenSuppressionSlice(input *[]string) []any {
	result := make([]any, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}
