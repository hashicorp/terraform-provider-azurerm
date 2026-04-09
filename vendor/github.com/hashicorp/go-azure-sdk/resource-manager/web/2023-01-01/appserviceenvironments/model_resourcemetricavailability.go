package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceMetricAvailability struct {
	Retention *string `json:"retention,omitempty"`
	TimeGrain *string `json:"timeGrain,omitempty"`
}
