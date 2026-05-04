package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionsAlwaysReadyConfig struct {
	InstanceCount *int64  `json:"instanceCount,omitempty"`
	Name          *string `json:"name,omitempty"`
}
