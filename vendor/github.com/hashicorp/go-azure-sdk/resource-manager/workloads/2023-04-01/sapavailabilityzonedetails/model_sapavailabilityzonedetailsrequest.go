package sapavailabilityzonedetails

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPAvailabilityZoneDetailsRequest struct {
	AppLocation  string          `json:"appLocation"`
	DatabaseType SAPDatabaseType `json:"databaseType"`
	SapProduct   SAPProductType  `json:"sapProduct"`
}
