package serveradministrators

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerAdministratorResource struct {
	Id         *string                        `json:"id,omitempty"`
	Name       *string                        `json:"name,omitempty"`
	Properties *ServerAdministratorProperties `json:"properties,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
