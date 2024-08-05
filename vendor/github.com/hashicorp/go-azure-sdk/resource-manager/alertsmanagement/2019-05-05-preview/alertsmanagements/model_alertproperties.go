package alertsmanagements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertProperties struct {
	Context      *interface{} `json:"context,omitempty"`
	EgressConfig *interface{} `json:"egressConfig,omitempty"`
	Essentials   *Essentials  `json:"essentials,omitempty"`
}
