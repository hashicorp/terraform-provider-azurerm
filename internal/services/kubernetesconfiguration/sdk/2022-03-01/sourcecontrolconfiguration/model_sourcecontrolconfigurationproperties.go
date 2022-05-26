package sourcecontrolconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceControlConfigurationProperties struct {
	ComplianceStatus               *ComplianceStatus       `json:"complianceStatus,omitempty"`
	ConfigurationProtectedSettings *map[string]string      `json:"configurationProtectedSettings,omitempty"`
	EnableHelmOperator             *bool                   `json:"enableHelmOperator,omitempty"`
	HelmOperatorProperties         *HelmOperatorProperties `json:"helmOperatorProperties,omitempty"`
	OperatorInstanceName           *string                 `json:"operatorInstanceName,omitempty"`
	OperatorNamespace              *string                 `json:"operatorNamespace,omitempty"`
	OperatorParams                 *string                 `json:"operatorParams,omitempty"`
	OperatorScope                  *OperatorScopeType      `json:"operatorScope,omitempty"`
	OperatorType                   *OperatorType           `json:"operatorType,omitempty"`
	ProvisioningState              *ProvisioningStateType  `json:"provisioningState,omitempty"`
	RepositoryPublicKey            *string                 `json:"repositoryPublicKey,omitempty"`
	RepositoryUrl                  *string                 `json:"repositoryUrl,omitempty"`
	SshKnownHostsContents          *string                 `json:"sshKnownHostsContents,omitempty"`
}
