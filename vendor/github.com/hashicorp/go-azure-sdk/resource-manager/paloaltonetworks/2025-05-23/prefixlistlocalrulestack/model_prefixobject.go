package prefixlistlocalrulestack

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrefixObject struct {
	AuditComment      *string            `json:"auditComment,omitempty"`
	Description       *string            `json:"description,omitempty"`
	Etag              *string            `json:"etag,omitempty"`
	PrefixList        []string           `json:"prefixList"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
