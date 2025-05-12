package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnosticDataProviderMetadata struct {
	PropertyBag  *[]DiagnosticDataProviderMetadataPropertyBagInlined `json:"propertyBag,omitempty"`
	ProviderName *string                                             `json:"providerName,omitempty"`
}
