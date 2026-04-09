package eventhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CaptureDescription struct {
	Destination       *Destination                `json:"destination,omitempty"`
	Enabled           *bool                       `json:"enabled,omitempty"`
	Encoding          *EncodingCaptureDescription `json:"encoding,omitempty"`
	IntervalInSeconds *int64                      `json:"intervalInSeconds,omitempty"`
	SizeLimitInBytes  *int64                      `json:"sizeLimitInBytes,omitempty"`
	SkipEmptyArchives *bool                       `json:"skipEmptyArchives,omitempty"`
}
