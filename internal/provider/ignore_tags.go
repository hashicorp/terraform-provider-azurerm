// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	rmtags "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// schemaIgnoreTags returns the schema for the provider-level `ignore_tags` block,
// which lets a user globally ignore specific tag keys (and key prefixes) across
// every resource and data source - mirroring the AWS provider's feature.
func schemaIgnoreTags() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Globally ignore tag changes for the specified keys and key prefixes across all resources and data sources.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"keys": {
					Type:         schema.TypeSet,
					Optional:     true,
					Elem:         &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringIsNotEmpty},
					AtLeastOneOf: []string{"ignore_tags.0.keys", "ignore_tags.0.key_prefixes"},
					Description:  "A set of tag keys, matched exactly (case-sensitive), to ignore across all resources and data sources.",
				},
				"key_prefixes": {
					Type:         schema.TypeSet,
					Optional:     true,
					Elem:         &schema.Schema{Type: schema.TypeString, ValidateFunc: validation.StringIsNotEmpty},
					AtLeastOneOf: []string{"ignore_tags.0.keys", "ignore_tags.0.key_prefixes"},
					Description:  "A set of tag key prefixes, matched case-sensitively, to ignore across all resources and data sources.",
				},
			},
		},
	}
}

// expandIgnoreTags parses the `ignore_tags` provider block into an IgnoreConfig.
// It returns nil when the block is absent or specifies no keys or key prefixes,
// which preserves the default behavior of not ignoring any tags.
func expandIgnoreTags(input []interface{}) *rmtags.IgnoreConfig {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})

	var keys, keyPrefixes []string
	if v, ok := raw["keys"].(*schema.Set); ok {
		for _, k := range v.List() {
			keys = append(keys, k.(string))
		}
	}
	if v, ok := raw["key_prefixes"].(*schema.Set); ok {
		for _, k := range v.List() {
			keyPrefixes = append(keyPrefixes, k.(string))
		}
	}

	if len(keys) == 0 && len(keyPrefixes) == 0 {
		return nil
	}

	return &rmtags.IgnoreConfig{
		Keys:        keys,
		KeyPrefixes: keyPrefixes,
	}
}
