package expressroutecrossconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCrossConnection struct {
	Etag       *string                                `json:"etag,omitempty"`
	Id         *string                                `json:"id,omitempty"`
	Location   *string                                `json:"location,omitempty"`
	Name       *string                                `json:"name,omitempty"`
	Properties *ExpressRouteCrossConnectionProperties `json:"properties,omitempty"`
	Tags       *map[string]string                     `json:"tags,omitempty"`
	Type       *string                                `json:"type,omitempty"`
}
