package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretsProfile struct {
	KeyVaultResourceId string             `json:"keyVaultResourceId"`
	Secrets            *[]SecretReference `json:"secrets,omitempty"`
}
