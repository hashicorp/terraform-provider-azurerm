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

type IngressConfigurationModel struct {
	WorkloadProfileName             string `tfschema:"workload_profile_name"`
	WorkloadProfileType             string `tfschema:"workload_profile_type"`
	MinimumNodeCount                int64  `tfschema:"minimum_node_count"`
	MaximumNodeCount                int64  `tfschema:"maximum_node_count"`
	TerminationGracePeriodMinutes   int64  `tfschema:"termination_grace_period_minutes"`
	RequestIdleTimeout              int64  `tfschema:"request_idle_timeout"`
	HeaderCountLimit                int64  `tfschema:"header_count_limit"`
}

func IngressConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:        pluginsdk.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "An `ingress_configuration` block as defined below. Configures Premium Ingress with a dedicated workload profile for ingress proxies.",
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"workload_profile_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      "premium-ingress",
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name of the dedicated workload profile for premium ingress. Defaults to `premium-ingress`.",
				},

				"workload_profile_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "D4",
					ValidateFunc: validation.StringInSlice([]string{
						string(WorkloadProfileSkuD4),
						string(WorkloadProfileSkuD8),
						string(WorkloadProfileSkuD16),
						string(WorkloadProfileSkuD32),
					}, false),
					Description: "The workload profile type for the dedicated ingress profile. Possible values are `D4`, `D8`, `D16`, and `D32`. Defaults to `D4`.",
				},

				"minimum_node_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      2,
					ValidateFunc: validation.IntAtLeast(2),
					Description:  "The minimum number of nodes for the dedicated ingress workload profile. Must be at least `2`. Defaults to `2`.",
				},

				"maximum_node_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      10,
					ValidateFunc: validation.IntAtLeast(2),
					Description:  "The maximum number of nodes for the dedicated ingress workload profile. Defaults to `10`.",
				},

				"termination_grace_period_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 60),
					Description:  "The termination grace period in minutes. The provider converts this to seconds for the API. If omitted, the backend default is used.",
				},

				"request_idle_timeout": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(4, 30),
					Description:  "The request idle timeout in minutes. If omitted, the backend default is used.",
				},

				"header_count_limit": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(1),
					Description:  "The maximum number of HTTP headers allowed per request. If omitted, the backend default is used.",
				},
			},
		},
	}
}

func IngressConfigurationSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:        pluginsdk.TypeList,
		Computed:    true,
		Description: "An `ingress_configuration` block as defined below.",
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"workload_profile_name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of the dedicated workload profile for premium ingress.",
				},
				"workload_profile_type": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The workload profile type for the dedicated ingress profile.",
				},
				"minimum_node_count": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The minimum number of nodes for the dedicated ingress workload profile.",
				},
				"maximum_node_count": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The maximum number of nodes for the dedicated ingress workload profile.",
				},
				"termination_grace_period_minutes": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The termination grace period in minutes.",
				},
				"request_idle_timeout": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The request idle timeout in minutes.",
				},
				"header_count_limit": {
					Type:        pluginsdk.TypeInt,
					Computed:    true,
					Description: "The maximum number of HTTP headers allowed per request.",
				},
			},
		},
	}
}

func ExpandIngressConfiguration(input []IngressConfigurationModel) (*managedenvironments.IngressConfiguration, *managedenvironments.WorkloadProfile) {
	if len(input) == 0 {
		return nil, nil
	}

	v := input[0]

	ingressConfig := &managedenvironments.IngressConfiguration{
		WorkloadProfileName: pointer.To(v.WorkloadProfileName),
	}

	// Only set these if explicitly configured; passing nil lets the backend use its own defaults
	if v.TerminationGracePeriodMinutes != 0 {
		ingressConfig.TerminationGracePeriodSeconds = pointer.To(v.TerminationGracePeriodMinutes * 60)
	}
	if v.RequestIdleTimeout != 0 {
		ingressConfig.RequestIdleTimeout = pointer.To(v.RequestIdleTimeout)
	}
	if v.HeaderCountLimit != 0 {
		ingressConfig.HeaderCountLimit = pointer.To(v.HeaderCountLimit)
	}

	workloadProfile := &managedenvironments.WorkloadProfile{
		Name:                v.WorkloadProfileName,
		WorkloadProfileType: v.WorkloadProfileType,
		MinimumCount:        pointer.To(v.MinimumNodeCount),
		MaximumCount:        pointer.To(v.MaximumNodeCount),
	}

	return ingressConfig, workloadProfile
}

func FlattenIngressConfiguration(ingressConfig *managedenvironments.IngressConfiguration, workloadProfiles *[]managedenvironments.WorkloadProfile) []IngressConfigurationModel {
	if ingressConfig == nil || ingressConfig.WorkloadProfileName == nil || *ingressConfig.WorkloadProfileName == "" {
		return []IngressConfigurationModel{}
	}

	result := IngressConfigurationModel{
		WorkloadProfileName:           pointer.From(ingressConfig.WorkloadProfileName),
		TerminationGracePeriodMinutes: pointer.From(ingressConfig.TerminationGracePeriodSeconds) / 60,
		RequestIdleTimeout:            pointer.From(ingressConfig.RequestIdleTimeout),
		HeaderCountLimit:              pointer.From(ingressConfig.HeaderCountLimit),
	}

	if workloadProfiles != nil && result.WorkloadProfileName != "" {
		for _, wp := range *workloadProfiles {
			if wp.Name == result.WorkloadProfileName {
				result.WorkloadProfileType = wp.WorkloadProfileType
				result.MinimumNodeCount = pointer.From(wp.MinimumCount)
				result.MaximumNodeCount = pointer.From(wp.MaximumCount)
				break
			}
		}
	}

	return []IngressConfigurationModel{result}
}

func MergeIngressWorkloadProfile(profiles *[]managedenvironments.WorkloadProfile, ingressProfile *managedenvironments.WorkloadProfile) *[]managedenvironments.WorkloadProfile {
	if ingressProfile == nil {
		return profiles
	}

	result := make([]managedenvironments.WorkloadProfile, 0)
	if profiles != nil {
		for _, p := range *profiles {
			if p.Name != ingressProfile.Name {
				result = append(result, p)
			}
		}
	}
	result = append(result, *ingressProfile)
	return &result
}

func FilterIngressWorkloadProfile(profiles []WorkloadProfileModel, ingressProfileName string) []WorkloadProfileModel {
	if ingressProfileName == "" {
		return profiles
	}

	result := make([]WorkloadProfileModel, 0)
	for _, p := range profiles {
		if p.Name != ingressProfileName {
			result = append(result, p)
		}
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
