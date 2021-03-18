package tags

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func Flatten(tagMap map[string]*string) map[string]interface{} {
	// If tagsMap is nil, len(tagsMap) will be 0.
	output := make(map[string]interface{}, len(tagMap))

	for i, v := range tagMap {
		if v == nil {
			continue
		}

		output[i] = *v
	}

	return output
}

func FlattenAndSet(d *schema.ResourceData, tagMap map[string]*string) error {
	flattened := Flatten(tagMap)
	if err := d.Set("tags", flattened); err != nil {
		return fmt.Errorf("Error setting `tags`: %s", err)
	}

	return nil
}

func IgnoreTags(tagMap map[string]*string, ignored features.IgnoreTagsFeatures) map[string]*string {
	ignoredKeys := ignored.Keys
	ignoredKeyPrefixes := ignored.KeyPrefixes
	for _, ignorePrefix := range ignoredKeyPrefixes {
		for key := range tagMap {
			if strings.HasPrefix(key, ignorePrefix) {
				delete(tagMap, key)
			}
		}
	}
	for _, ignore := range ignoredKeys {
		if _, ok := tagMap[ignore]; ok {
			delete(tagMap, ignore)
		}
	}

	return tagMap
}
