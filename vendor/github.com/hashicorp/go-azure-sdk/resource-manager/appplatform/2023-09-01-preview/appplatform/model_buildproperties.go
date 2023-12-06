package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BuildProperties struct {
	AgentPool            *string                 `json:"agentPool,omitempty"`
	Apms                 *[]ApmReference         `json:"apms,omitempty"`
	Builder              *string                 `json:"builder,omitempty"`
	Certificates         *[]CertificateReference `json:"certificates,omitempty"`
	Env                  *map[string]string      `json:"env,omitempty"`
	ProvisioningState    *BuildProvisioningState `json:"provisioningState,omitempty"`
	RelativePath         *string                 `json:"relativePath,omitempty"`
	ResourceRequests     *BuildResourceRequests  `json:"resourceRequests,omitempty"`
	TriggeredBuildResult *TriggeredBuildResult   `json:"triggeredBuildResult,omitempty"`
}
