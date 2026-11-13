package indexes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalyzedTokenInfo struct {
	EndOffset   int64  `json:"endOffset"`
	Position    int64  `json:"position"`
	StartOffset int64  `json:"startOffset"`
	Token       string `json:"token"`
}
