package mongorbacs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoRoleDefinitionType int64

const (
	MongoRoleDefinitionTypeOne  MongoRoleDefinitionType = 1
	MongoRoleDefinitionTypeZero MongoRoleDefinitionType = 0
)

func PossibleValuesForMongoRoleDefinitionType() []int64 {
	return []int64{
		int64(MongoRoleDefinitionTypeOne),
		int64(MongoRoleDefinitionTypeZero),
	}
}
