package elasticsanskus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuInformationList struct {
	NextLink *string           `json:"nextLink,omitempty"`
	Value    *[]SkuInformation `json:"value,omitempty"`
}
