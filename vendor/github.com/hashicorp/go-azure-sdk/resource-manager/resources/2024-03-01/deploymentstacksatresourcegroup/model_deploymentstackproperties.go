package deploymentstacksatresourcegroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentStackProperties struct {
	ActionOnUnmanage          ActionOnUnmanage                  `json:"actionOnUnmanage"`
	BypassStackOutOfSyncError *bool                             `json:"bypassStackOutOfSyncError,omitempty"`
	CorrelationId             *string                           `json:"correlationId,omitempty"`
	DebugSetting              *DeploymentStacksDebugSetting     `json:"debugSetting,omitempty"`
	DeletedResources          *[]ResourceReference              `json:"deletedResources,omitempty"`
	DenySettings              DenySettings                      `json:"denySettings"`
	DeploymentId              *string                           `json:"deploymentId,omitempty"`
	DeploymentScope           *string                           `json:"deploymentScope,omitempty"`
	Description               *string                           `json:"description,omitempty"`
	DetachedResources         *[]ResourceReference              `json:"detachedResources,omitempty"`
	Duration                  *string                           `json:"duration,omitempty"`
	Error                     *ErrorDetail                      `json:"error,omitempty"`
	FailedResources           *[]ResourceReferenceExtended      `json:"failedResources,omitempty"`
	Outputs                   *map[string]interface{}           `json:"outputs,omitempty"`
	Parameters                *map[string]DeploymentParameter   `json:"parameters,omitempty"`
	ParametersLink            *DeploymentStacksParametersLink   `json:"parametersLink,omitempty"`
	ProvisioningState         *DeploymentStackProvisioningState `json:"provisioningState,omitempty"`
	Resources                 *[]ManagedResourceReference       `json:"resources,omitempty"`
	Template                  *map[string]interface{}           `json:"template,omitempty"`
	TemplateLink              *DeploymentStacksTemplateLink     `json:"templateLink,omitempty"`
}
