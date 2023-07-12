package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationProperties struct {
	ApplicationDefinitionId *string                              `json:"applicationDefinitionId,omitempty"`
	Artifacts               *[]ApplicationArtifact               `json:"artifacts,omitempty"`
	Authorizations          *[]ApplicationAuthorization          `json:"authorizations,omitempty"`
	BillingDetails          *ApplicationBillingDetailsDefinition `json:"billingDetails,omitempty"`
	CreatedBy               *ApplicationClientDetails            `json:"createdBy,omitempty"`
	CustomerSupport         *ApplicationPackageContact           `json:"customerSupport,omitempty"`
	JitAccessPolicy         *ApplicationJitAccessPolicy          `json:"jitAccessPolicy,omitempty"`
	ManagedResourceGroupId  *string                              `json:"managedResourceGroupId,omitempty"`
	ManagementMode          *ApplicationManagementMode           `json:"managementMode,omitempty"`
	Outputs                 *interface{}                         `json:"outputs,omitempty"`
	Parameters              *interface{}                         `json:"parameters,omitempty"`
	ProvisioningState       *ProvisioningState                   `json:"provisioningState,omitempty"`
	PublisherTenantId       *string                              `json:"publisherTenantId,omitempty"`
	SupportUrls             *ApplicationPackageSupportUrls       `json:"supportUrls,omitempty"`
	UpdatedBy               *ApplicationClientDetails            `json:"updatedBy,omitempty"`
}
