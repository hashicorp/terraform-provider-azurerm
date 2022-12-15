package cognitiveservicesaccounts

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountModel struct {
	BaseModel     *DeploymentModel       `json:"baseModel,omitempty"`
	CallRateLimit *CallRateLimit         `json:"callRateLimit,omitempty"`
	Capabilities  *map[string]string     `json:"capabilities,omitempty"`
	Deprecation   *ModelDeprecationInfo  `json:"deprecation,omitempty"`
	Format        *string                `json:"format,omitempty"`
	MaxCapacity   *int64                 `json:"maxCapacity,omitempty"`
	Name          *string                `json:"name,omitempty"`
	SystemData    *systemdata.SystemData `json:"systemData,omitempty"`
	Version       *string                `json:"version,omitempty"`
}
