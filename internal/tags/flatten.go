// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"fmt"

	rmtags "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func Flatten(tagMap map[string]*string) map[string]any {
	// scrub any provider-level ignored tag keys before they reach state
	tagMap = rmtags.Ignore().ApplyMap(tagMap)

	// If tagsMap is nil, len(tagsMap) will be 0.
	output := make(map[string]any, len(tagMap))

	for i, v := range tagMap {
		if v == nil {
			continue
		}

		output[i] = *v
	}

	return output
}

func FlattenAndSet(d *pluginsdk.ResourceData, tagMap map[string]*string) error {
	if err := d.Set("tags", Flatten(tagMap)); err != nil {
		return fmt.Errorf("setting `tags`: %s", err)
	}

	return nil
}
