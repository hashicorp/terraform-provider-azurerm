package activitylogs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPRequestInfo struct {
	ClientIPAddress *string `json:"clientIpAddress,omitempty"`
	ClientRequestId *string `json:"clientRequestId,omitempty"`
	Method          *string `json:"method,omitempty"`
	Uri             *string `json:"uri,omitempty"`
}
