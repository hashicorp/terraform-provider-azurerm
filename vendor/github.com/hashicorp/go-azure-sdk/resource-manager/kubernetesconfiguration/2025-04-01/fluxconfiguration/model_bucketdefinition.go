package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BucketDefinition struct {
	AccessKey             *string `json:"accessKey,omitempty"`
	BucketName            *string `json:"bucketName,omitempty"`
	Insecure              *bool   `json:"insecure,omitempty"`
	LocalAuthRef          *string `json:"localAuthRef,omitempty"`
	SyncIntervalInSeconds *int64  `json:"syncIntervalInSeconds,omitempty"`
	TimeoutInSeconds      *int64  `json:"timeoutInSeconds,omitempty"`
	Url                   *string `json:"url,omitempty"`
}
