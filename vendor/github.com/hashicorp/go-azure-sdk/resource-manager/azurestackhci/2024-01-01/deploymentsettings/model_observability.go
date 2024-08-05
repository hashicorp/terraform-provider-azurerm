package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Observability struct {
	EpisodicDataUpload  *bool `json:"episodicDataUpload,omitempty"`
	EuLocation          *bool `json:"euLocation,omitempty"`
	StreamingDataClient *bool `json:"streamingDataClient,omitempty"`
}
