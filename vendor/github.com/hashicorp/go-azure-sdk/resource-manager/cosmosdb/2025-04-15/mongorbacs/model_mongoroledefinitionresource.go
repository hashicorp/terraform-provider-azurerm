package mongorbacs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoRoleDefinitionResource struct {
	DatabaseName *string                  `json:"databaseName,omitempty"`
	Privileges   *[]Privilege             `json:"privileges,omitempty"`
	RoleName     *string                  `json:"roleName,omitempty"`
	Roles        *[]Role                  `json:"roles,omitempty"`
	Type         *MongoRoleDefinitionType `json:"type,omitempty"`
}
