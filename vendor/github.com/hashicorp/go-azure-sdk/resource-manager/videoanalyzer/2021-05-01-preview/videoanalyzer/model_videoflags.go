package videoanalyzer

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VideoFlags struct {
	CanStream   bool `json:"canStream"`
	HasData     bool `json:"hasData"`
	IsRecording bool `json:"isRecording"`
}
