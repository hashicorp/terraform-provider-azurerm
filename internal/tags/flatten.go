// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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

func FlattenAndSet(d *pluginsdk.ResourceData, tagMap map[string]*string) error {
	flattened := Flatten(tagMap)
	if err := d.Set("tags", flattened); err != nil {
		return fmt.Errorf("setting `tags`: %s", err)
	}

	return nil
}
