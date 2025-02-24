package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SasPolicy struct {
	ExpirationAction    ExpirationAction `json:"expirationAction"`
	SasExpirationPeriod string           `json:"sasExpirationPeriod"`
}
