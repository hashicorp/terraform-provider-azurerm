package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ItemLevelRestoreCriteria = KubernetesClusterVaultTierRestoreCriteria{}

type KubernetesClusterVaultTierRestoreCriteria struct {
	ConflictPolicy               *ExistingResourcePolicy      `json:"conflictPolicy,omitempty"`
	ExcludedNamespaces           *[]string                    `json:"excludedNamespaces,omitempty"`
	ExcludedResourceTypes        *[]string                    `json:"excludedResourceTypes,omitempty"`
	IncludeClusterScopeResources bool                         `json:"includeClusterScopeResources"`
	IncludedNamespaces           *[]string                    `json:"includedNamespaces,omitempty"`
	IncludedResourceTypes        *[]string                    `json:"includedResourceTypes,omitempty"`
	LabelSelectors               *[]string                    `json:"labelSelectors,omitempty"`
	NamespaceMappings            *map[string]string           `json:"namespaceMappings,omitempty"`
	PersistentVolumeRestoreMode  *PersistentVolumeRestoreMode `json:"persistentVolumeRestoreMode,omitempty"`
	ResourceModifierReference    *NamespacedNameResource      `json:"resourceModifierReference,omitempty"`
	RestoreHookReferences        *[]NamespacedNameResource    `json:"restoreHookReferences,omitempty"`
	StagingResourceGroupId       *string                      `json:"stagingResourceGroupId,omitempty"`
	StagingStorageAccountId      *string                      `json:"stagingStorageAccountId,omitempty"`

	// Fields inherited from ItemLevelRestoreCriteria

	ObjectType string `json:"objectType"`
}

func (s KubernetesClusterVaultTierRestoreCriteria) ItemLevelRestoreCriteria() BaseItemLevelRestoreCriteriaImpl {
	return BaseItemLevelRestoreCriteriaImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = KubernetesClusterVaultTierRestoreCriteria{}

func (s KubernetesClusterVaultTierRestoreCriteria) MarshalJSON() ([]byte, error) {
	type wrapper KubernetesClusterVaultTierRestoreCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KubernetesClusterVaultTierRestoreCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KubernetesClusterVaultTierRestoreCriteria: %+v", err)
	}

	decoded["objectType"] = "KubernetesClusterVaultTierRestoreCriteria"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KubernetesClusterVaultTierRestoreCriteria: %+v", err)
	}

	return encoded, nil
}
