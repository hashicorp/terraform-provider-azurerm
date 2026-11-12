package cloudvmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortRange struct {
	Max int64 `json:"max"`
	Min int64 `json:"min"`
}
