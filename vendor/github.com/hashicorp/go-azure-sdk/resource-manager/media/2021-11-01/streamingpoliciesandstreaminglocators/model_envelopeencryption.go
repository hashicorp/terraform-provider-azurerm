package streamingpoliciesandstreaminglocators

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvelopeEncryption struct {
	ClearTracks                     *[]TrackSelection           `json:"clearTracks,omitempty"`
	ContentKeys                     *StreamingPolicyContentKeys `json:"contentKeys,omitempty"`
	CustomKeyAcquisitionUrlTemplate *string                     `json:"customKeyAcquisitionUrlTemplate,omitempty"`
	EnabledProtocols                *EnabledProtocols           `json:"enabledProtocols,omitempty"`
}
