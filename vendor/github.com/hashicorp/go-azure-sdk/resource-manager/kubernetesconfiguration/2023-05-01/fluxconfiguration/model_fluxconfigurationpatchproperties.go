package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluxConfigurationPatchProperties struct {
	AzureBlob                      *AzureBlobPatchDefinition                `json:"azureBlob,omitempty"`
	Bucket                         *BucketPatchDefinition                   `json:"bucket,omitempty"`
	ConfigurationProtectedSettings *map[string]string                       `json:"configurationProtectedSettings,omitempty"`
	GitRepository                  *GitRepositoryPatchDefinition            `json:"gitRepository,omitempty"`
	Kustomizations                 *map[string]KustomizationPatchDefinition `json:"kustomizations,omitempty"`
	SourceKind                     *SourceKindType                          `json:"sourceKind,omitempty"`
	Suspend                        *bool                                    `json:"suspend,omitempty"`
}
