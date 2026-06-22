package cognitiveservicesaccounts

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountModel struct {
	BaseModel            *DeploymentModel       `json:"baseModel,omitempty"`
	CallRateLimit        *CallRateLimit         `json:"callRateLimit,omitempty"`
	Capabilities         *map[string]string     `json:"capabilities,omitempty"`
	Deprecation          *ModelDeprecationInfo  `json:"deprecation,omitempty"`
	FinetuneCapabilities *map[string]string     `json:"finetuneCapabilities,omitempty"`
	Format               *string                `json:"format,omitempty"`
	IsDefaultVersion     *bool                  `json:"isDefaultVersion,omitempty"`
	LifecycleStatus      *ModelLifecycleStatus  `json:"lifecycleStatus,omitempty"`
	MaxCapacity          *int64                 `json:"maxCapacity,omitempty"`
	ModelCatalogAssetId  *string                `json:"modelCatalogAssetId,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	Publisher            *string                `json:"publisher,omitempty"`
	ReplacementConfig    *ReplacementConfig     `json:"replacementConfig,omitempty"`
	Skus                 *[]ModelSku            `json:"skus,omitempty"`
	Source               *string                `json:"source,omitempty"`
	SourceAccount        *string                `json:"sourceAccount,omitempty"`
	SystemData           *systemdata.SystemData `json:"systemData,omitempty"`
	Version              *string                `json:"version,omitempty"`
}
