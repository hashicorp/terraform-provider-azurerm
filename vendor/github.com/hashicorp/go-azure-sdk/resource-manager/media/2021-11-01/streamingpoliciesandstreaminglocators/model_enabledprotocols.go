package streamingpoliciesandstreaminglocators

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnabledProtocols struct {
	Dash            bool `json:"dash"`
	Download        bool `json:"download"`
	Hls             bool `json:"hls"`
	SmoothStreaming bool `json:"smoothStreaming"`
}
