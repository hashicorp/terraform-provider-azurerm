package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdxDestination struct {
	DatabaseName *string `json:"databaseName,omitempty"`
	IngestionUri *string `json:"ingestionUri,omitempty"`
	Name         *string `json:"name,omitempty"`
	ResourceId   *string `json:"resourceId,omitempty"`
}
