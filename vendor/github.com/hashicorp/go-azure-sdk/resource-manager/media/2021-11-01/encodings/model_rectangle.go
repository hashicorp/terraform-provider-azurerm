package encodings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Rectangle struct {
	Height *string `json:"height,omitempty"`
	Left   *string `json:"left,omitempty"`
	Top    *string `json:"top,omitempty"`
	Width  *string `json:"width,omitempty"`
}
