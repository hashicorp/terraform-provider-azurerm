package applicationdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationDefinitionProperties struct {
	Artifacts          *[]ApplicationDefinitionArtifact           `json:"artifacts,omitempty"`
	Authorizations     *[]ApplicationAuthorization                `json:"authorizations,omitempty"`
	CreateUiDefinition *interface{}                               `json:"createUiDefinition,omitempty"`
	DeploymentPolicy   *ApplicationDeploymentPolicy               `json:"deploymentPolicy,omitempty"`
	Description        *string                                    `json:"description,omitempty"`
	DisplayName        *string                                    `json:"displayName,omitempty"`
	IsEnabled          *bool                                      `json:"isEnabled,omitempty"`
	LockLevel          ApplicationLockLevel                       `json:"lockLevel"`
	LockingPolicy      *ApplicationPackageLockingPolicyDefinition `json:"lockingPolicy,omitempty"`
	MainTemplate       *interface{}                               `json:"mainTemplate,omitempty"`
	ManagementPolicy   *ApplicationManagementPolicy               `json:"managementPolicy,omitempty"`
	NotificationPolicy *ApplicationNotificationPolicy             `json:"notificationPolicy,omitempty"`
	PackageFileUri     *string                                    `json:"packageFileUri,omitempty"`
	Policies           *[]ApplicationPolicy                       `json:"policies,omitempty"`
	StorageAccountId   *string                                    `json:"storageAccountId,omitempty"`
}
