// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package service_plan

type ServicePlanModel struct {
	Name                        string            `tfschema:"name"`
	ResourceGroup               string            `tfschema:"resource_group_name"`
	Location                    string            `tfschema:"location"`
	Kind                        string            `tfschema:"kind"`
	OSType                      OSType            `tfschema:"os_type"`
	Sku                         string            `tfschema:"sku_name"`
	AppServiceEnvironmentId     string            `tfschema:"app_service_environment_id"`
	PerSiteScaling              bool              `tfschema:"per_site_scaling_enabled"`
	Reserved                    bool              `tfschema:"reserved"`
	WorkerCount                 int64             `tfschema:"worker_count"`
	PremiumPlanAutoScaleEnabled bool              `tfschema:"premium_plan_auto_scale_enabled"`
	MaximumElasticWorkerCount   int64             `tfschema:"maximum_elastic_worker_count"`
	ZoneBalancing               bool              `tfschema:"zone_balancing_enabled"`
	Tags                        map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &ServicePlanModel{}
}
