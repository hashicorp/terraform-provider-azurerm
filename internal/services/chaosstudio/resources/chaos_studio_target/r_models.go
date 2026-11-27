// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_target

type ChaosStudioTargetResourceSchema struct {
	Location         string `tfschema:"location"`
	TargetResourceId string `tfschema:"target_resource_id"`
	TargetType       string `tfschema:"target_type"`
}

func (r Resource) ModelObject() any {
	return &ChaosStudioTargetResourceSchema{}
}
