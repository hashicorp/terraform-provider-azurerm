package nodereports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscReportError struct {
	ErrorCode    *string `json:"errorCode,omitempty"`
	ErrorDetails *string `json:"errorDetails,omitempty"`
	ErrorMessage *string `json:"errorMessage,omitempty"`
	ErrorSource  *string `json:"errorSource,omitempty"`
	Locale       *string `json:"locale,omitempty"`
	ResourceId   *string `json:"resourceId,omitempty"`
}
