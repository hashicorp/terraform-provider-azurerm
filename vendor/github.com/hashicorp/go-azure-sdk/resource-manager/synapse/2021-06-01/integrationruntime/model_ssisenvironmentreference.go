package integrationruntime

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SsisEnvironmentReference struct {
	EnvironmentFolderName *string `json:"environmentFolderName,omitempty"`
	EnvironmentName       *string `json:"environmentName,omitempty"`
	Id                    *int64  `json:"id,omitempty"`
	ReferenceType         *string `json:"referenceType,omitempty"`
}
