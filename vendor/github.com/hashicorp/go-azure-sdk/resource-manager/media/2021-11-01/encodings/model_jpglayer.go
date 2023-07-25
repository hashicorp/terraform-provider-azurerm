package encodings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JpgLayer struct {
	Height  *string `json:"height,omitempty"`
	Label   *string `json:"label,omitempty"`
	Quality *int64  `json:"quality,omitempty"`
	Width   *string `json:"width,omitempty"`
}
