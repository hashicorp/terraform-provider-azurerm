package environmenttypes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserRoleAssignment struct {
	Roles *map[string]EnvironmentRole `json:"roles,omitempty"`
}
