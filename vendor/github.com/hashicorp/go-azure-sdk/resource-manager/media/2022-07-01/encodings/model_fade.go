package encodings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Fade struct {
	Duration  string  `json:"duration"`
	FadeColor string  `json:"fadeColor"`
	Start     *string `json:"start,omitempty"`
}
