// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"regexp"
	"strings"
	"sync"
)

// PropertyExceptions contains properties and the rules they should skip for.
// Format: property pattern -> rule IDs
var PropertyExceptions = map[string]map[string]struct{}{
	// Timezone properties - some values cannot be fully enumerated
	"all.timezone":     {"S011": {}},
	"all.time_zone":    {"S011": {}},
	"all.time_zone_id": {"S011": {}},

	// Key Vault references - these reference external resources
	"all.key_vault_secret_id":      {"S011": {}},
	"all.key_vault_key_id":         {"S011": {}},
	"all.key_vault_certificate_id": {"S011": {}},

	"/azurerm_orchestrated_virtual_machine_scale_set.os_profile.(windows|linux)_configuration.secret.certificate": {"S009": {}}, // Not declared in schema, but certificate expects a `store` param
	"/azurerm_eventgrid_event_subscription.advanced_filter.*":                                                     {"S009": {}}, // nested properties are generalized and not marked as BLOCk in doc
	"/azurerm_eventgrid_system_topic_event_subscription.advanced_filter.*":                                        {"S009": {}}, // nested properties are generalized and not marked as BLOCk in doc
}

var (
	skipPathsByRule  map[string]map[string]struct{} // ruleID -> property paths
	skipRegexpByRule map[string][]*regexp.Regexp    // ruleID -> regex patterns
	initSkipOnce     sync.Once
)

func initPropertySkipConfig() {
	initSkipOnce.Do(func() {
		skipPathsByRule = make(map[string]map[string]struct{})
		skipRegexpByRule = make(map[string][]*regexp.Regexp)

		for property, rules := range PropertyExceptions {
			for ruleID := range rules {
				if strings.HasPrefix(property, "/") {
					// It's a regex pattern
					if skipRegexpByRule[ruleID] == nil {
						skipRegexpByRule[ruleID] = make([]*regexp.Regexp, 0)
					}
					skipRegexpByRule[ruleID] = append(skipRegexpByRule[ruleID], regexp.MustCompile(property[1:]))
				} else {
					// It's a simple path
					if skipPathsByRule[ruleID] == nil {
						skipPathsByRule[ruleID] = make(map[string]struct{})
					}
					skipPathsByRule[ruleID][property] = struct{}{}
				}
			}
		}
	})
}

// SkipProp checks if a property should be skipped for a specific rule
// ruleID is the rule identifier (e.g., "S007", "S011")
// rt is the resource type (e.g., "azurerm_storage_account")
// prop is the property path (e.g., "blob_properties.cors_rule")
func SkipProp(ruleID, rt, prop string) bool {
	initPropertySkipConfig()

	skipPaths := skipPathsByRule[ruleID]
	skipRegexps := skipRegexpByRule[ruleID]

	// Check if entire resource is skipped
	if _, ok := skipPaths[rt]; ok {
		return true
	}

	// Check if specific property path is skipped
	fullPath := rt + "." + prop
	if _, ok := skipPaths[fullPath]; ok {
		return true
	}

	// Check if property name (last part) matches an "all.xxx" pattern
	propName := prop
	if idx := strings.LastIndex(prop, "."); idx >= 0 {
		propName = prop[idx+1:]
	}
	if _, ok := skipPaths["all."+propName]; ok {
		return true
	}

	// Check against regex patterns
	for _, reg := range skipRegexps {
		if reg.MatchString(fullPath) {
			return true
		}
	}

	return false
}
