package integrationruntimedisableinteractivequery

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PackageStore struct {
	Name                      string          `json:"name"`
	PackageStoreLinkedService EntityReference `json:"packageStoreLinkedService"`
}
