package tags

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

func FlattenAndSetTags(d *schema.ResourceData, tags interface{}) error {
	tagMap := tags.(map[string]interface{})
	currentTagMap := make(map[string]*string)

	for k, v := range tagMap {
		currentTagMap[k] = utils.String(v.(string))
	}

	flattened := Flatten(currentTagMap)
	if err := d.Set("tags", flattened); err != nil {
		return fmt.Errorf("Error setting `tags`: %s", err)
	}

	return nil
}
