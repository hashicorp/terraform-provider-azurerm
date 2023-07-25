package iscsitargets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Acl struct {
	InitiatorIqn string   `json:"initiatorIqn"`
	MappedLuns   []string `json:"mappedLuns"`
}
