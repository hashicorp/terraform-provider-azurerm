package secrets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretCreateOrUpdateParameters struct {
	Properties SecretProperties   `json:"properties"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
