package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceEndpoint struct {
	Locations *[]string `json:"locations,omitempty"`
	Service   string    `json:"service"`
}
