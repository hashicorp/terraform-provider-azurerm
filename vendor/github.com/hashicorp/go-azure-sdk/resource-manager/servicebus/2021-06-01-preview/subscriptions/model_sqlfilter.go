package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlFilter struct {
	CompatibilityLevel    *int64  `json:"compatibilityLevel,omitempty"`
	RequiresPreprocessing *bool   `json:"requiresPreprocessing,omitempty"`
	SqlExpression         *string `json:"sqlExpression,omitempty"`
}
