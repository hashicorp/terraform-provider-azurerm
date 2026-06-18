package mongorbacs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoUserDefinitionResource struct {
	CustomData   *string `json:"customData,omitempty"`
	DatabaseName *string `json:"databaseName,omitempty"`
	Mechanisms   *string `json:"mechanisms,omitempty"`
	Password     *string `json:"password,omitempty"`
	Roles        *[]Role `json:"roles,omitempty"`
	UserName     *string `json:"userName,omitempty"`
}
