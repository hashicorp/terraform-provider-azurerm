package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BackupDatasourceParameters = KubernetesClusterBackupDatasourceParameters{}

type KubernetesClusterBackupDatasourceParameters struct {
	BackupHookReferences         *[]NamespacedNameResource `json:"backupHookReferences,omitempty"`
	ExcludedNamespaces           *[]string                 `json:"excludedNamespaces,omitempty"`
	ExcludedResourceTypes        *[]string                 `json:"excludedResourceTypes,omitempty"`
	IncludeClusterScopeResources bool                      `json:"includeClusterScopeResources"`
	IncludedNamespaces           *[]string                 `json:"includedNamespaces,omitempty"`
	IncludedResourceTypes        *[]string                 `json:"includedResourceTypes,omitempty"`
	LabelSelectors               *[]string                 `json:"labelSelectors,omitempty"`
	SnapshotVolumes              bool                      `json:"snapshotVolumes"`

	// Fields inherited from BackupDatasourceParameters

	ObjectType string `json:"objectType"`
}

func (s KubernetesClusterBackupDatasourceParameters) BackupDatasourceParameters() BaseBackupDatasourceParametersImpl {
	return BaseBackupDatasourceParametersImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = KubernetesClusterBackupDatasourceParameters{}

func (s KubernetesClusterBackupDatasourceParameters) MarshalJSON() ([]byte, error) {
	type wrapper KubernetesClusterBackupDatasourceParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KubernetesClusterBackupDatasourceParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KubernetesClusterBackupDatasourceParameters: %+v", err)
	}

	decoded["objectType"] = "KubernetesClusterBackupDatasourceParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KubernetesClusterBackupDatasourceParameters: %+v", err)
	}

	return encoded, nil
}
