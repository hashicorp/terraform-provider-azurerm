package streamingendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AkamaiAccessControl struct {
	AkamaiSignatureHeaderAuthenticationKeyList *[]AkamaiSignatureHeaderAuthenticationKey `json:"akamaiSignatureHeaderAuthenticationKeyList,omitempty"`
}
