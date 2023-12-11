package environmenttypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentRole struct {
	Description *string `json:"description,omitempty"`
	RoleName    *string `json:"roleName,omitempty"`
}
