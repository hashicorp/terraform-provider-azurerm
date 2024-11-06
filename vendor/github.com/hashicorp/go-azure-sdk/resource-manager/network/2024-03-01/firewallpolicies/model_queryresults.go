package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryResults struct {
	MatchingRecordsCount *int64               `json:"matchingRecordsCount,omitempty"`
	Signatures           *[]SingleQueryResult `json:"signatures,omitempty"`
}
