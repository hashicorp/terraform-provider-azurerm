package metadata

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetadataDependencies struct {
	ContentId *string                 `json:"contentId,omitempty"`
	Criteria  *[]MetadataDependencies `json:"criteria,omitempty"`
	Kind      *Kind                   `json:"kind,omitempty"`
	Name      *string                 `json:"name,omitempty"`
	Operator  *Operator               `json:"operator,omitempty"`
	Version   *string                 `json:"version,omitempty"`
}
