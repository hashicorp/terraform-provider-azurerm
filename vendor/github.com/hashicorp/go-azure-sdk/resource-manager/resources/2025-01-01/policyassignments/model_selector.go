package policyassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Selector struct {
	In    *[]string     `json:"in,omitempty"`
	Kind  *SelectorKind `json:"kind,omitempty"`
	NotIn *[]string     `json:"notIn,omitempty"`
}
