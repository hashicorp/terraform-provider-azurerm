package healthbots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthBotKeysResponse struct {
	Secrets *[]HealthBotKey `json:"secrets,omitempty"`
}
