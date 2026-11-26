package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckSkuAvailabilityParameter struct {
	Kind string   `json:"kind"`
	Skus []string `json:"skus"`
	Type string   `json:"type"`
}
