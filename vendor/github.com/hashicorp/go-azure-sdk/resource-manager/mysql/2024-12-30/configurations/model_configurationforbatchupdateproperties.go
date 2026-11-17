package configurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationForBatchUpdateProperties struct {
	Source *string `json:"source,omitempty"`
	Value  *string `json:"value,omitempty"`
}
