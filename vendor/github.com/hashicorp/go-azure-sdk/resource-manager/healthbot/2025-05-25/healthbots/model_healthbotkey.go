package healthbots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthBotKey struct {
	KeyName *string `json:"keyName,omitempty"`
	Value   *string `json:"value,omitempty"`
}
