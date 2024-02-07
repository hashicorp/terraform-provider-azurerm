package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkResource struct {
	Id         string                        `json:"id"`
	Name       string                        `json:"name"`
	Properties PrivateLinkResourceProperties `json:"properties"`
	Type       string                        `json:"type"`
}
