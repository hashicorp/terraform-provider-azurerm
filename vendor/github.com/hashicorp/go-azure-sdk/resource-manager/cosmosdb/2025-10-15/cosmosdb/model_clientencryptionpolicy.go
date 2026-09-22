package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClientEncryptionPolicy struct {
	IncludedPaths       []ClientEncryptionIncludedPath `json:"includedPaths"`
	PolicyFormatVersion int64                          `json:"policyFormatVersion"`
}
