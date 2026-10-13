package expressroutecircuits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLinkFailoverRoute struct {
	NextHop         *string `json:"nextHop,omitempty"`
	PrimaryASPath   *string `json:"primaryASPath,omitempty"`
	Route           *string `json:"route,omitempty"`
	SecondaryASPath *string `json:"secondaryASPath,omitempty"`
}
