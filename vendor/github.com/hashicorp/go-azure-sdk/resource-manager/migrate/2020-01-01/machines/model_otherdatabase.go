package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OtherDatabase struct {
	DatabaseType *string `json:"databaseType,omitempty"`
	Instance     *string `json:"instance,omitempty"`
	Version      *string `json:"version,omitempty"`
}
