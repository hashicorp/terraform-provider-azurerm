package fhirservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceOciArtifactEntry struct {
	Digest      *string `json:"digest,omitempty"`
	ImageName   *string `json:"imageName,omitempty"`
	LoginServer *string `json:"loginServer,omitempty"`
}
