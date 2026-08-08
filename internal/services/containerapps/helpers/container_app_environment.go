// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WorkloadProfileSku string

// NOTE: the Workload Profile SKUs aren't defined in the Swagger definition so we define them here
const (
	WorkloadProfileSkuConsumption            WorkloadProfileSku = "Consumption"
	WorkloadProfileSkuConsumptionGpuNc24A100 WorkloadProfileSku = "Consumption-GPU-NC24-A100"
	WorkloadProfileSkuConsumptionGpuNc8AsT4  WorkloadProfileSku = "Consumption-GPU-NC8as-T4"
	WorkloadProfileSkuD4                     WorkloadProfileSku = "D4"
	WorkloadProfileSkuD8                     WorkloadProfileSku = "D8"
	WorkloadProfileSkuD16                    WorkloadProfileSku = "D16"
	WorkloadProfileSkuD32                    WorkloadProfileSku = "D32"
	WorkloadProfileSkuE4                     WorkloadProfileSku = "E4"
	WorkloadProfileSkuE8                     WorkloadProfileSku = "E8"
	WorkloadProfileSkuE16                    WorkloadProfileSku = "E16"
	WorkloadProfileSkuE32                    WorkloadProfileSku = "E32"
	WorkloadProfileSkuNc24A100               WorkloadProfileSku = "NC24-A100"
	WorkloadProfileSkuNc48A100               WorkloadProfileSku = "NC48-A100"
	WorkloadProfileSkuNc96A100               WorkloadProfileSku = "NC96-A100"
)

func PossibleValuesForWorkloadProfileSku() []string {
	return []string{
		string(WorkloadProfileSkuConsumption),
		string(WorkloadProfileSkuConsumptionGpuNc24A100),
		string(WorkloadProfileSkuConsumptionGpuNc8AsT4),
		string(WorkloadProfileSkuD4),
		string(WorkloadProfileSkuD8),
		string(WorkloadProfileSkuD16),
		string(WorkloadProfileSkuD32),
		string(WorkloadProfileSkuE4),
		string(WorkloadProfileSkuE8),
		string(WorkloadProfileSkuE16),
		string(WorkloadProfileSkuE32),
		string(WorkloadProfileSkuNc24A100),
		string(WorkloadProfileSkuNc48A100),
		string(WorkloadProfileSkuNc96A100),
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
		Type:                  pluginsdk.TypeSet,
		Optional:              true,
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, _, _ string, d *pluginsdk.ResourceData) bool {
			o, n := d.GetChange("workload_profile")

			oldProfiles := o.(*pluginsdk.Set)
			newProfiles := n.(*pluginsdk.Set)

			return OneAdditionalConsumptionProfileReturnedByAPI(oldProfiles, newProfiles)
		},
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

func FlattenWorkloadProfiles(input *[]managedenvironments.WorkloadProfile) []WorkloadProfileModel {
	if input == nil || len(*input) == 0 {
		return []WorkloadProfileModel{}
	}
	result := make([]WorkloadProfileModel, 0)

	for _, v := range *input {
		result = append(result, WorkloadProfileModel{
			Name:                v.Name,
			MaximumCount:        pointer.From(v.MaximumCount),
			MinimumCount:        pointer.From(v.MinimumCount),
			WorkloadProfileType: v.WorkloadProfileType,
		})
	}

	return result
}

func OneAdditionalConsumptionProfileReturnedByAPI(returnedProfiles, definedProfiles *pluginsdk.Set) bool {
	// if 1 more profile is returned by the API than is defined, then check if it is a consumption profile
	if returnedProfiles.Len() == definedProfiles.Len()+1 {
		// check if we have defined a consumption profile
		for _, v := range definedProfiles.List() {
			profile := v.(map[string]interface{})
			if profile["workload_profile_type"].(string) == string(WorkloadProfileSkuConsumption) {
				return false
			}
		}

		// now that we know there are no consumption profiles defined in the config, check if the API returned a consumption profile
		for _, v := range returnedProfiles.List() {
			profile := v.(map[string]interface{})
			if profile["workload_profile_type"].(string) == string(WorkloadProfileSkuConsumption) {
				return true
			}
		}
	}
	return false
}
