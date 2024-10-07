package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperatingSystem struct {
	OsName    *string `json:"osName,omitempty"`
	OsType    *string `json:"osType,omitempty"`
	OsVersion *string `json:"osVersion,omitempty"`
}
