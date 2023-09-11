package assessmentsmetadata

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityAssessmentMetadataPartnerData struct {
	PartnerName string  `json:"partnerName"`
	ProductName *string `json:"productName,omitempty"`
	Secret      string  `json:"secret"`
}
