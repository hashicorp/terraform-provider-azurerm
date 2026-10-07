package account

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessKeys struct {
	AtlasKafkaPrimaryEndpoint   *string `json:"atlasKafkaPrimaryEndpoint,omitempty"`
	AtlasKafkaSecondaryEndpoint *string `json:"atlasKafkaSecondaryEndpoint,omitempty"`
}
