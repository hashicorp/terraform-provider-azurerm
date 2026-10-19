package communicationservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkNotificationHubParameters struct {
	ConnectionString string `json:"connectionString"`
	ResourceId       string `json:"resourceId"`
}
