package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SSISExecutionCredential struct {
	Domain   string       `json:"domain"`
	Password SecureString `json:"password"`
	UserName string       `json:"userName"`
}
