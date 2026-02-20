package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Processor struct {
	Name          *string `json:"name,omitempty"`
	NumberOfCores *int64  `json:"numberOfCores,omitempty"`
}
