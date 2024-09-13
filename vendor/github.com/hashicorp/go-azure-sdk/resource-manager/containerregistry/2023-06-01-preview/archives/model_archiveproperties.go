package archives

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArchiveProperties struct {
	PackageSource            *ArchivePackageSourceProperties `json:"packageSource,omitempty"`
	ProvisioningState        *ProvisioningState              `json:"provisioningState,omitempty"`
	PublishedVersion         *string                         `json:"publishedVersion,omitempty"`
	RepositoryEndpoint       *string                         `json:"repositoryEndpoint,omitempty"`
	RepositoryEndpointPrefix *string                         `json:"repositoryEndpointPrefix,omitempty"`
}
