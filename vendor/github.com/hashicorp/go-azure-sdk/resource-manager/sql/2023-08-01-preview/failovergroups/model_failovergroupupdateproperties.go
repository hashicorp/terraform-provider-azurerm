package failovergroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailoverGroupUpdateProperties struct {
	Databases         *[]string                            `json:"databases,omitempty"`
	PartnerServers    *[]PartnerInfo                       `json:"partnerServers,omitempty"`
	ReadOnlyEndpoint  *FailoverGroupReadOnlyEndpoint       `json:"readOnlyEndpoint,omitempty"`
	ReadWriteEndpoint *FailoverGroupReadWriteEndpoint      `json:"readWriteEndpoint,omitempty"`
	SecondaryType     *FailoverGroupDatabasesSecondaryType `json:"secondaryType,omitempty"`
}
