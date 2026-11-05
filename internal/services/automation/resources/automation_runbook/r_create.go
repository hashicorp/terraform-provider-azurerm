// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_runbook

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
)

func expandContentLink(inputs []any) *runbook.ContentLink {
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}

	input := inputs[0].(map[string]any)
	uri := input["uri"].(string)
	version := input["version"].(string)
	hashes := input["hash"].([]any)

	if len(hashes) > 0 {
		hash := hashes[0].(map[string]any)
		hashValue := hash["value"].(string)
		hashAlgorithm := hash["algorithm"].(string)

		return &runbook.ContentLink{
			Uri:     &uri,
			Version: &version,
			ContentHash: &runbook.ContentHash{
				Algorithm: hashAlgorithm,
				Value:     hashValue,
			},
		}
	}

	return &runbook.ContentLink{
		Uri:     &uri,
		Version: &version,
	}
}

func expandDraft(inputs []any) *runbook.RunbookDraft {
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}

	input := inputs[0].(map[string]any)
	var res runbook.RunbookDraft

	res.DraftContentLink = expandContentLink(input["content_link"].([]any))
	res.InEdit = pointer.To(input["edit_mode_enabled"].(bool))
	parameter := map[string]runbook.RunbookParameter{}

	for _, iparam := range input["parameters"].([]any) {
		param := iparam.(map[string]any)
		key := param["key"].(string)
		parameter[key] = runbook.RunbookParameter{
			Type:         pointer.To(param["type"].(string)),
			IsMandatory:  pointer.To(param["mandatory"].(bool)),
			Position:     pointer.To(int64(param["position"].(int))),
			DefaultValue: pointer.To(param["default_value"].(string)),
		}
	}
	res.Parameters = &parameter

	typesInput := input["output_types"].([]any)
	types := make([]string, 0, len(typesInput))
	for _, v := range typesInput {
		types = append(types, v.(string))
	}

	if len(types) > 0 {
		res.OutputTypes = &types
	}

	return &res
}
