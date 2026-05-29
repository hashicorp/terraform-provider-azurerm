package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Expression struct {
	Type  ExpressionType `json:"type"`
	Value string         `json:"value"`
}
