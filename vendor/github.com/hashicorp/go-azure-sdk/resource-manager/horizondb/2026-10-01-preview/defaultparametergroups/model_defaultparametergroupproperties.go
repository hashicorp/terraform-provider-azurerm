package defaultparametergroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DefaultParameterGroupProperties struct {
	Description       *string               `json:"description,omitempty"`
	Parameters        []ParameterProperties `json:"parameters"`
	PostgreSqlVersion PostgreSqlVersion     `json:"postgreSqlVersion"`
	ProvisioningState *ProvisioningState    `json:"provisioningState,omitempty"`
}
