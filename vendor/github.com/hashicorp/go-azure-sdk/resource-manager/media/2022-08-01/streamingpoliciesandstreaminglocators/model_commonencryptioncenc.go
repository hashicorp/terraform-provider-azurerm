package streamingpoliciesandstreaminglocators

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommonEncryptionCenc struct {
	ClearKeyEncryptionConfiguration *ClearKeyEncryptionConfiguration `json:"clearKeyEncryptionConfiguration,omitempty"`
	ClearTracks                     *[]TrackSelection                `json:"clearTracks,omitempty"`
	ContentKeys                     *StreamingPolicyContentKeys      `json:"contentKeys,omitempty"`
	Drm                             *CencDrmConfiguration            `json:"drm,omitempty"`
	EnabledProtocols                *EnabledProtocols                `json:"enabledProtocols,omitempty"`
}
