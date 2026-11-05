// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_deployment

type cognitiveDeploymentModel struct {
	Name                     string                 `tfschema:"name"`
	CognitiveAccountId       string                 `tfschema:"cognitive_account_id"`
	DynamicThrottlingEnabled bool                   `tfschema:"dynamic_throttling_enabled"`
	Model                    []DeploymentModelModel `tfschema:"model"`
	RaiPolicyName            string                 `tfschema:"rai_policy_name"`
	Sku                      []DeploymentSkuModel   `tfschema:"sku"`
	SpilloverDeploymentName  string                 `tfschema:"spillover_deployment_name"`
	VersionUpgradeOption     string                 `tfschema:"version_upgrade_option"`
}

type DeploymentModelModel struct {
	Format  string `tfschema:"format"`
	Name    string `tfschema:"name"`
	Version string `tfschema:"version"`
}

type DeploymentSkuModel struct {
	Name     string `tfschema:"name"`
	Tier     string `tfschema:"tier"`
	Size     string `tfschema:"size"`
	Family   string `tfschema:"family"`
	Capacity int64  `tfschema:"capacity"`
}

func (r Resource) ModelObject() any {
	return &cognitiveDeploymentModel{}
}
