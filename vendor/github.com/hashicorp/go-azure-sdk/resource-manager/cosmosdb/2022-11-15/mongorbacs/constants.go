package mongorbacs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoRoleDefinitionType int64

const (
	MongoRoleDefinitionTypeBuiltInRole MongoRoleDefinitionType = 0
	MongoRoleDefinitionTypeCustomRole  MongoRoleDefinitionType = 1
)

func PossibleValuesForMongoRoleDefinitionType() []int64 {
	return []int64{
		int64(MongoRoleDefinitionTypeBuiltInRole),
		int64(MongoRoleDefinitionTypeCustomRole),
	}
}
