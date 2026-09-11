// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CustomDiffComputedSet bypasses a bug/quirk in Terraform Plugin SDKv2 where
// elements of a Computed TypeSet have their inner diffs completely suppressed
// if their identity hash matches.
//
// Removals from the set are treated natively (i.e. as an in-place update to the Set).
func CustomDiffComputedSet(setKey string, hashFunc schema.SchemaSetFunc, forceNewProps []string, inPlaceProps []string) CustomizeDiffFunc {
	return customDiffComputedSet(setKey, false, hashFunc, forceNewProps, inPlaceProps)
}

// CustomDiffComputedSetCannotRemove is identical to CustomDiffComputedSet, but
// explicitly forces a replacement of the entire resource if an element is removed
// from the Set.
func CustomDiffComputedSetCannotRemove(setKey string, hashFunc schema.SchemaSetFunc, forceNewProps []string, inPlaceProps []string) CustomizeDiffFunc {
	return customDiffComputedSet(setKey, true, hashFunc, forceNewProps, inPlaceProps)
}

func customDiffComputedSet(setKey string, forceNewOnRemoval bool, hashFunc schema.SchemaSetFunc, forceNewProps []string, inPlaceProps []string) CustomizeDiffFunc {
	return func(ctx context.Context, diff *ResourceDiff, meta interface{}) error {
		oldRaw, _ := diff.GetChange(setKey)
		var oldSet *schema.Set
		if oldRaw != nil {
			oldSet = oldRaw.(*schema.Set)
		}

		oldMap := make(map[int]map[string]interface{})
		if oldSet != nil {
			for _, v := range oldSet.List() {
				if m, ok := v.(map[string]interface{}); ok {
					hash := hashFunc(m)
					oldMap[hash] = m
				}
			}
		}

		rawConfig := diff.GetRawConfig()
		if rawConfig.IsNull() || !rawConfig.IsKnown() {
			return nil
		}
		rawConfigMap := rawConfig.AsValueMap()
		rawSetCty, ok := rawConfigMap[setKey]
		if !ok || rawSetCty.IsNull() || !rawSetCty.IsKnown() {
			return nil
		}

		var forceNew bool
		newMap := make(map[int]bool)

		for _, itemCty := range rawSetCty.AsValueSet().Values() {
			itemMap := itemCty.AsValueMap()

			// Convert cty.Value to map[string]interface{} for hashFunc
			inputMap := make(map[string]interface{})
			for k, v := range itemMap {
				if v.IsNull() || !v.IsKnown() {
					continue
				}
				switch v.Type() {
				case cty.String:
					inputMap[k] = v.AsString()
				case cty.Number:
					if f, _ := v.AsBigFloat().Float64(); true {
						// SDKv2 parses numbers as float64
						inputMap[k] = f
					}
				case cty.Bool:
					inputMap[k] = v.True()
				}
			}

			hash := hashFunc(inputMap)
			newMap[hash] = true

			oldInput, exists := oldMap[hash]
			if !exists {
				// Hash doesn't match anything in old state, SDKv2 will treat this
				// natively as an addition, so we don't need to patch it.
				continue
			}

			updateField := func(prop string, triggersForceNew bool) {
				var newPropStr string
				if val, ok := itemMap[prop]; ok && !val.IsNull() && val.IsKnown() {
					if val.Type() == cty.String {
						newPropStr = val.AsString()
					} else {
						// Fallback string conversion for primitive non-strings
						newPropStr = fmt.Sprintf("%v", inputMap[prop])
					}
				}

				var oldPropStr string
				if oldVal, ok := oldInput[prop]; ok && oldVal != nil {
					oldPropStr = fmt.Sprintf("%v", oldVal)
				}

				if newPropStr != oldPropStr {
					if triggersForceNew {
						forceNew = true
					} else {
						// For in-place updates, manually patch the flatmap so it survives.
						key := fmt.Sprintf("%s.%d.%s", setKey, hash, prop)
						_ = diff.SetNew(key, newPropStr)
					}
				}
			}

			for _, prop := range forceNewProps {
				updateField(prop, true)
			}
			for _, prop := range inPlaceProps {
				updateField(prop, false)
			}
		}

		if forceNewOnRemoval {
			for hash := range oldMap {
				if !newMap[hash] {
					forceNew = true
				}
			}
		}

		if forceNew {
			if err := diff.SetNewComputed(setKey); err != nil {
				return err
			}
			if err := diff.ForceNew(setKey); err != nil {
				return err
			}
		}

		return nil
	}
}
