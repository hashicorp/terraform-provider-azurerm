package autonomousdatabases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GenerateAutonomousDatabaseWalletDetails struct {
	GenerateType *GenerateType `json:"generateType,omitempty"`
	IsRegional   *bool         `json:"isRegional,omitempty"`
	Password     string        `json:"password"`
}
