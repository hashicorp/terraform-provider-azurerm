// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_solution

type SolutionResourceModel struct {
	SolutionName        string              `tfschema:"solution_name"`
	WorkspaceName       string              `tfschema:"workspace_name"`
	WorkspaceResourceId string              `tfschema:"workspace_resource_id"`
	Location            string              `tfschema:"location"`
	ResourceGroupName   string              `tfschema:"resource_group_name"`
	SolutionPlan        []SolutionPlanModel `tfschema:"plan"`
	Tags                map[string]string   `tfschema:"tags"`
}

type SolutionPlanModel struct {
	Name          string `tfschema:"name"`
	Publisher     string `tfschema:"publisher"`
	PromotionCode string `tfschema:"promotion_code"`
	Product       string `tfschema:"product"`
}

func (s Resource) ModelObject() any {
	return &SolutionResourceModel{}
}
