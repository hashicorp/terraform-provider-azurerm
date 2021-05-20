package advisor

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/advisor/mgmt/2020-01-01/advisor"
	"github.com/gofrs/uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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
						string(advisor.HighAvailability),
						string(advisor.Security),
						string(advisor.Performance),
						string(advisor.Cost),
						string(advisor.OperationalExcellence),
					}, true),
				},
			},

			"filter_by_resource_groups": azure.SchemaResourceGroupNameSetOptional(),

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

	filterList := make([]string, 0)
	if categories := expandAzureRmAdvisorRecommendationsMapString("Category", d.Get("filter_by_category").(*pluginsdk.Set).List()); categories != "" {
		filterList = append(filterList, categories)
	}
	if resGroups := expandAzureRmAdvisorRecommendationsMapString("ResourceGroup", d.Get("filter_by_resource_groups").(*pluginsdk.Set).List()); resGroups != "" {
		filterList = append(filterList, resGroups)
	}

	var recommends []advisor.ResourceRecommendationBase
	for recommendationIterator, err := client.ListComplete(ctx, strings.Join(filterList, " and "), nil, ""); recommendationIterator.NotDone(); err = recommendationIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading Advisor Recommendation List: %+v", err)
		}

		if recommendationIterator.Value().Name == nil || *recommendationIterator.Value().Name == "" {
			return fmt.Errorf("advisor Recommendation Name was nil or empty")
		}

		recommends = append(recommends, recommendationIterator.Value())
	}

	if err := d.Set("recommendations", flattenAzureRmAdvisorRecommendations(recommends)); err != nil {
		return fmt.Errorf("setting `recommendations`: %+v", err)
	}

	d.SetId(fmt.Sprintf("avdisor/recommendations/%s", time.Now().UTC().String()))

	return nil
}

func flattenAzureRmAdvisorRecommendations(recommends []advisor.ResourceRecommendationBase) []interface{} {
	result := make([]interface{}, 0)

	if len(recommends) == 0 {
		return result
	}

	for _, v := range recommends {
		var category, description, impact, recTypeId, resourceName, resourceType, updatedTime string
		var suppressionIds []interface{}
		if v.Category != "" {
			category = string(v.Category)
		}

		if v.ShortDescription != nil && v.ShortDescription.Problem != nil {
			description = *v.ShortDescription.Problem
		}

		if v.Impact != "" {
			impact = string(v.Impact)
		}

		if v.RecommendationTypeID != nil {
			recTypeId = *v.RecommendationTypeID
		}

		if v.ImpactedValue != nil {
			resourceName = *v.ImpactedValue
		}

		if v.ImpactedField != nil {
			resourceType = *v.ImpactedField
		}

		if v.SuppressionIds != nil {
			suppressionIds = flattenSuppressionSlice(v.SuppressionIds)
		}
		if v.LastUpdated != nil && !v.LastUpdated.IsZero() {
			updatedTime = v.LastUpdated.Format(time.RFC3339)
		}

		result = append(result, map[string]interface{}{
			"category":               category,
			"description":            description,
			"impact":                 impact,
			"recommendation_name":    *v.Name,
			"recommendation_type_id": recTypeId,
			"resource_name":          resourceName,
			"resource_type":          resourceType,
			"suppression_names":      suppressionIds,
			"updated_time":           updatedTime,
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

func flattenSuppressionSlice(input *[]uuid.UUID) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item.String())
		}
	}
	return result
}
