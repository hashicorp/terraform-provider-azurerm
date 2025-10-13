package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubDataSource struct {
	ConsumerGroup *string `json:"consumerGroup,omitempty"`
	Name          *string `json:"name,omitempty"`
	Stream        *string `json:"stream,omitempty"`
}
