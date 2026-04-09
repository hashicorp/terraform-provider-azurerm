package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdifactProcessingSettings struct {
	CreateEmptyXmlTagsForTrailingSeparators bool `json:"createEmptyXmlTagsForTrailingSeparators"`
	MaskSecurityInfo                        bool `json:"maskSecurityInfo"`
	PreserveInterchange                     bool `json:"preserveInterchange"`
	SuspendInterchangeOnError               bool `json:"suspendInterchangeOnError"`
	UseDotAsDecimalSeparator                bool `json:"useDotAsDecimalSeparator"`
}
