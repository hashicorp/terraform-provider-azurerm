package automation

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func tagValueToString(v interface{}) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil
	case int:
		return fmt.Sprintf("%d", value), nil
	default:
		return "", fmt.Errorf("unknown tag type %T in tag value", value)
	}
}

func expandTags(tagsMap map[string]interface{}) map[string]string {
	output := make(map[string]string, len(tagsMap))

	for i, v := range tagsMap {
		// Validate should have ignored this error already
		value, _ := tagValueToString(v)
		output[i] = value
	}

	return output
}

func flattenTags(tagMap map[string]string) map[string]interface{} {
	// If tagsMap is nil, len(tagsMap) will be 0.
	output := make(map[string]interface{}, len(tagMap))

	for i, v := range tagMap {
		output[i] = v
	}

	return output
}

func flattenAndSetTags(d *pluginsdk.ResourceData, tagMap map[string]string) error {
	flattened := flattenTags(tagMap)
	if err := d.Set("tags", flattened); err != nil {
		return fmt.Errorf("setting `tags`: %s", err)
	}

	return nil
}
