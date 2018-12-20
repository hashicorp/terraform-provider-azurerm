package azure

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func ExpandApplicationInsightsAPIKeyLinkedProperties(v *schema.Set, subscriptionID, resGroup, appInsightsName string) *[]string {
	if v == nil {
		return nil
	}

	baseLinkedProperty := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/microsoft.insights/components/%s", subscriptionID, resGroup, appInsightsName)

	result := make([]string, v.Len())
	for i, prop := range v.List() {
		result[i] = fmt.Sprintf("%s/%s", baseLinkedProperty, prop)
	}
	return &result
}

func FlattenApplicationInsightsAPIKeyLinkedProperties(props *[]string) *[]string {
	if props == nil {
		return nil
	}

	result := make([]string, len(*props))
	for i, prop := range *props {
		elems := strings.Split(prop, "/")
		result[i] = elems[len(elems)-1]
	}
	return &result
}
