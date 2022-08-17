package backupinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetDetails struct {
	FilePrefix                string                    `json:"filePrefix"`
	RestoreTargetLocationType RestoreTargetLocationType `json:"restoreTargetLocationType"`
	Url                       string                    `json:"url"`
}
