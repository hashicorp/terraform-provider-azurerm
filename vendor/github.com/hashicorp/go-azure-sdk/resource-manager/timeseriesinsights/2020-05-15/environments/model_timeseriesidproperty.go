package environments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TimeSeriesIdProperty struct {
	Name *string       `json:"name,omitempty"`
	Type *PropertyType `json:"type,omitempty"`
}
