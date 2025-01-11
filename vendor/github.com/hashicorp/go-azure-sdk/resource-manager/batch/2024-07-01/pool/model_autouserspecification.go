package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoUserSpecification struct {
	ElevationLevel *ElevationLevel `json:"elevationLevel,omitempty"`
	Scope          *AutoUserScope  `json:"scope,omitempty"`
}
