package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportSource struct {
	Credentials *ImportSourceCredentials `json:"credentials,omitempty"`
	RegistryUri *string                  `json:"registryUri,omitempty"`
	ResourceId  *string                  `json:"resourceId,omitempty"`
	SourceImage string                   `json:"sourceImage"`
}
