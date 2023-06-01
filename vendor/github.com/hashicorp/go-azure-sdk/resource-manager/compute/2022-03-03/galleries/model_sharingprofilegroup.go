package galleries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharingProfileGroup struct {
	Ids  *[]string                 `json:"ids,omitempty"`
	Type *SharingProfileGroupTypes `json:"type,omitempty"`
}
