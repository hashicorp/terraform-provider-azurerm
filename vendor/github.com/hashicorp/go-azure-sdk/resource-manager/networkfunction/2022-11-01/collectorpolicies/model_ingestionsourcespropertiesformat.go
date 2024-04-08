package collectorpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IngestionSourcesPropertiesFormat struct {
	ResourceId *string     `json:"resourceId,omitempty"`
	SourceType *SourceType `json:"sourceType,omitempty"`
}
