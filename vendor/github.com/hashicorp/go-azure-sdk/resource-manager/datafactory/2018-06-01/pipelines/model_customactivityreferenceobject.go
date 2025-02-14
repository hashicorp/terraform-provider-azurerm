package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomActivityReferenceObject struct {
	Datasets       *[]DatasetReference       `json:"datasets,omitempty"`
	LinkedServices *[]LinkedServiceReference `json:"linkedServices,omitempty"`
}
