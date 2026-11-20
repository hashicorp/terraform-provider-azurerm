package graphqueries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GraphQueryPropertiesUpdateParameters struct {
	Description *string `json:"description,omitempty"`
	Query       *string `json:"query,omitempty"`
}
