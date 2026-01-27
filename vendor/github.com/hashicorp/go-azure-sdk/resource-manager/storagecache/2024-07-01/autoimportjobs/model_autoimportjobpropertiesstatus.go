package autoimportjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoImportJobPropertiesStatus struct {
	AutoImportProgress *int64                `json:"autoImportProgress,omitempty"`
	State              *AutoImportStatusType `json:"state,omitempty"`
}
