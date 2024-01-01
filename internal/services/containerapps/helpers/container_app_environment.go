// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const consumption = "Consumption"

type WorkloadProfileModel struct {
	MaximumCount        int    `tfschema:"maximum_count"`
	MinimumCount        int    `tfschema:"minimum_count"`
	Name                string `tfschema:"name"`
	WorkloadProfileType string `tfschema:"workload_profile_type"`
}

func WorkloadProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeSet,
		Optional:     true,
		RequiredWith: []string{"infrastructure_resource_group_name"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"workload_profile_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						"D4",
						"D8",
						"D16",
						"D32",
						"E4",
						"E8",
						"E16",
						"E32",
					}, false),
				},

				"maximum_count": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},

				"minimum_count": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
			},
		},
	}
}

func ExpandWorkloadProfiles(input []WorkloadProfileModel) *[]managedenvironments.WorkloadProfile {
	if len(input) == 0 {
		return nil
	}

	result := make([]managedenvironments.WorkloadProfile, 0)

	for _, v := range input {
		r := managedenvironments.WorkloadProfile{
			Name:                v.Name,
			WorkloadProfileType: v.WorkloadProfileType,
		}

		if v.Name != consumption {
			r.MaximumCount = pointer.To(int64(v.MaximumCount))
			r.MinimumCount = pointer.To(int64(v.MinimumCount))
		}

		result = append(result, r)
	}

	result = append(result, managedenvironments.WorkloadProfile{
		Name:                consumption,
		WorkloadProfileType: consumption,
	})

	return &result
}

func FlattenWorkloadProfiles(input *[]managedenvironments.WorkloadProfile) []WorkloadProfileModel {
	if input == nil || len(*input) == 0 {
		return []WorkloadProfileModel{}
	}
	result := make([]WorkloadProfileModel, 0)

	for _, v := range *input {
		if strings.EqualFold(v.WorkloadProfileType, consumption) {
			continue
		}
		result = append(result, WorkloadProfileModel{
			Name:                v.Name,
			MaximumCount:        int(pointer.From(v.MaximumCount)),
			MinimumCount:        int(pointer.From(v.MinimumCount)),
			WorkloadProfileType: v.WorkloadProfileType,
		})
	}

	return result
}
