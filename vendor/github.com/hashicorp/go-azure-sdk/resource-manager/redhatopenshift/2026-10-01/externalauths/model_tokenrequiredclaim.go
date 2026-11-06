package externalauths

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TokenRequiredClaim struct {
	Claim         string `json:"claim"`
	RequiredValue string `json:"requiredValue"`
}
