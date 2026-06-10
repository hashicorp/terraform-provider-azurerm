package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NfsAccessPolicy struct {
	AccessRules []NfsAccessRule `json:"accessRules"`
	Name        string          `json:"name"`
}
