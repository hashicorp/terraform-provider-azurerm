package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessListServiceAccountsSuccessResponse struct {
	Data     *[]ServiceAccountRecord `json:"data,omitempty"`
	Kind     *string                 `json:"kind,omitempty"`
	Metadata *ConfluentListMetadata  `json:"metadata,omitempty"`
}
