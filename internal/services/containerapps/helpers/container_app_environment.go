// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WorkloadProfileSku string

// NOTE: the Workload Profile SKUs aren't defined in the Swagger definition so we define them here
const (
	WorkloadProfileSkuConsumption WorkloadProfileSku = "Consumption"
	WorkloadProfileSkuD4          WorkloadProfileSku = "D4"
	WorkloadProfileSkuD8          WorkloadProfileSku = "D8"
	WorkloadProfileSkuD16         WorkloadProfileSku = "D16"
	WorkloadProfileSkuD32         WorkloadProfileSku = "D32"
	WorkloadProfileSkuE4          WorkloadProfileSku = "E4"
	WorkloadProfileSkuE8          WorkloadProfileSku = "E8"
	WorkloadProfileSkuE16         WorkloadProfileSku = "E16"
	WorkloadProfileSkuE32         WorkloadProfileSku = "E32"
)

func PossibleValuesForWorkloadProfileSku() []string {
	return []string{
		string(WorkloadProfileSkuConsumption),
		string(WorkloadProfileSkuD4),
		string(WorkloadProfileSkuD8),
		string(WorkloadProfileSkuD16),
		string(WorkloadProfileSkuD32),
		string(WorkloadProfileSkuE4),
		string(WorkloadProfileSkuE8),
		string(WorkloadProfileSkuE16),
		string(WorkloadProfileSkuE32),
	}
}

type WorkloadProfileModel struct {
	MaximumCount        int64  `tfschema:"maximum_count"`
	MinimumCount        int64  `tfschema:"minimum_count"`
	Name                string `tfschema:"name"`
	WorkloadProfileType string `tfschema:"workload_profile_type"`
}

func WorkloadProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"workload_profile_type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(PossibleValuesForWorkloadProfileSku(), false),
				},

				"maximum_count": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},

				"minimum_count": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
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
			Name: v.Name,
		}

		if v.Name != string(WorkloadProfileSkuConsumption) {
			r.WorkloadProfileType = v.WorkloadProfileType
			r.MaximumCount = pointer.To(v.MaximumCount)
			r.MinimumCount = pointer.To(v.MinimumCount)
		} else {
			r.WorkloadProfileType = string(WorkloadProfileSkuConsumption)
		}

		result = append(result, r)
	}

	return &result
}

func FlattenWorkloadProfiles(input *[]managedenvironments.WorkloadProfile, consumptionDefined bool) []WorkloadProfileModel {
	if input == nil || len(*input) == 0 {
		return []WorkloadProfileModel{}
	}
	result := make([]WorkloadProfileModel, 0)

	for _, v := range *input {
		if strings.EqualFold(v.WorkloadProfileType, string(WorkloadProfileSkuConsumption)) && !consumptionDefined {
			continue
		}
		result = append(result, WorkloadProfileModel{
			Name:                v.Name,
			MaximumCount:        pointer.From(v.MaximumCount),
			MinimumCount:        pointer.From(v.MinimumCount),
			WorkloadProfileType: v.WorkloadProfileType,
		})
	}

	return result
}
