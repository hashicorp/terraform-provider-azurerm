package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationProperties struct {
	ApplicationState          *string                        `json:"applicationState,omitempty"`
	ApplicationType           *string                        `json:"applicationType,omitempty"`
	ComputeProfile            *ComputeProfile                `json:"computeProfile,omitempty"`
	CreatedDate               *string                        `json:"createdDate,omitempty"`
	Errors                    *[]Errors                      `json:"errors,omitempty"`
	HTTPSEndpoints            *[]ApplicationGetHTTPSEndpoint `json:"httpsEndpoints,omitempty"`
	InstallScriptActions      *[]RuntimeScriptAction         `json:"installScriptActions,omitempty"`
	MarketplaceIdentifier     *string                        `json:"marketplaceIdentifier,omitempty"`
	PrivateLinkConfigurations *[]PrivateLinkConfiguration    `json:"privateLinkConfigurations,omitempty"`
	ProvisioningState         *string                        `json:"provisioningState,omitempty"`
	SshEndpoints              *[]ApplicationGetEndpoint      `json:"sshEndpoints,omitempty"`
	UninstallScriptActions    *[]RuntimeScriptAction         `json:"uninstallScriptActions,omitempty"`
}
