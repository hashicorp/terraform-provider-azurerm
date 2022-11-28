package streamingpoliciesandstreaminglocators

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CbcsDrmConfiguration struct {
	FairPlay  *StreamingPolicyFairPlayConfiguration  `json:"fairPlay"`
	PlayReady *StreamingPolicyPlayReadyConfiguration `json:"playReady"`
	Widevine  *StreamingPolicyWidevineConfiguration  `json:"widevine"`
}
