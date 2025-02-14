package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SSISLogLocation struct {
	LogPath        string                        `json:"logPath"`
	Type           SsisLogLocationType           `json:"type"`
	TypeProperties SSISLogLocationTypeProperties `json:"typeProperties"`
}
