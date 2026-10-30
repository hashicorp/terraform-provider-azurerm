package apis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiProperties struct {
	Contacts              *[]Contact               `json:"contacts,omitempty"`
	CustomProperties      *interface{}             `json:"customProperties,omitempty"`
	Description           *string                  `json:"description,omitempty"`
	ExternalDocumentation *[]ExternalDocumentation `json:"externalDocumentation,omitempty"`
	Kind                  ApiKind                  `json:"kind"`
	License               *License                 `json:"license,omitempty"`
	LifecycleStage        *LifecycleStage          `json:"lifecycleStage,omitempty"`
	Summary               *string                  `json:"summary,omitempty"`
	TermsOfService        *TermsOfService          `json:"termsOfService,omitempty"`
	Title                 string                   `json:"title"`
}
