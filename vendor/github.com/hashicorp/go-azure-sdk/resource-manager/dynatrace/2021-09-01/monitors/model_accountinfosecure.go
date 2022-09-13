package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountInfoSecure struct {
	AccountId *string `json:"accountId,omitempty"`
	ApiKey    *string `json:"apiKey,omitempty"`
	RegionId  *string `json:"regionId,omitempty"`
}
