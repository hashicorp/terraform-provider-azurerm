package oraclesubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivationLinks struct {
	ExistingCloudAccountActivationLink *string `json:"existingCloudAccountActivationLink,omitempty"`
	NewCloudAccountActivationLink      *string `json:"newCloudAccountActivationLink,omitempty"`
}
