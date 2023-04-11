package streamingpoliciesandstreaminglocators

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingPolicyPlayReadyConfiguration struct {
	CustomLicenseAcquisitionUrlTemplate *string `json:"customLicenseAcquisitionUrlTemplate,omitempty"`
	PlayReadyCustomAttributes           *string `json:"playReadyCustomAttributes,omitempty"`
}
