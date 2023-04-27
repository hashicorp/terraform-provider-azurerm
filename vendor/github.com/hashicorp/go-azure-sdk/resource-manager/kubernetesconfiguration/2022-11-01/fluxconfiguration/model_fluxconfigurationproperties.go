package fluxconfiguration

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluxConfigurationProperties struct {
	AzureBlob                      *AzureBlobDefinition                `json:"azureBlob,omitempty"`
	Bucket                         *BucketDefinition                   `json:"bucket,omitempty"`
	ComplianceState                *FluxComplianceState                `json:"complianceState,omitempty"`
	ConfigurationProtectedSettings *map[string]string                  `json:"configurationProtectedSettings,omitempty"`
	ErrorMessage                   *string                             `json:"errorMessage,omitempty"`
	GitRepository                  *GitRepositoryDefinition            `json:"gitRepository,omitempty"`
	Kustomizations                 *map[string]KustomizationDefinition `json:"kustomizations,omitempty"`
	Namespace                      *string                             `json:"namespace,omitempty"`
	ProvisioningState              *ProvisioningState                  `json:"provisioningState,omitempty"`
	RepositoryPublicKey            *string                             `json:"repositoryPublicKey,omitempty"`
	Scope                          *ScopeType                          `json:"scope,omitempty"`
	SourceKind                     *SourceKindType                     `json:"sourceKind,omitempty"`
	SourceSyncedCommitId           *string                             `json:"sourceSyncedCommitId,omitempty"`
	SourceUpdatedAt                *string                             `json:"sourceUpdatedAt,omitempty"`
	StatusUpdatedAt                *string                             `json:"statusUpdatedAt,omitempty"`
	Statuses                       *[]ObjectStatusDefinition           `json:"statuses,omitempty"`
	Suspend                        *bool                               `json:"suspend,omitempty"`
}

func (o *FluxConfigurationProperties) GetSourceUpdatedAtAsTime() (*time.Time, error) {
	if o.SourceUpdatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SourceUpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *FluxConfigurationProperties) SetSourceUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SourceUpdatedAt = &formatted
}

func (o *FluxConfigurationProperties) GetStatusUpdatedAtAsTime() (*time.Time, error) {
	if o.StatusUpdatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StatusUpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *FluxConfigurationProperties) SetStatusUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StatusUpdatedAt = &formatted
}
