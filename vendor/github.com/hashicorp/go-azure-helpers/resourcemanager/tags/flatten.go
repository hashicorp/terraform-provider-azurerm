// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Flatten transforms the Tags specified via `input` into a map[string]interface{}
// for compatibility with the Schema.
func Flatten(input *map[string]string) map[string]interface{} {
	output := make(map[string]interface{})
	if input == nil {
		return output
	}

	for k, v := range *input {
		tagKey := k
		tagValue := v
		output[tagKey] = tagValue
	}

	return output
}

// FlattenAndSet first Flatten's the Tags and then sets the flattened value into
// the `tags` field in the State.
func FlattenAndSet(d *schema.ResourceData, input *map[string]string) error {
	tags := Flatten(input)

	if err := d.Set("tags", tags); err != nil {
		return fmt.Errorf("setting `tags`: %+v", err)
	}

	return nil
}
