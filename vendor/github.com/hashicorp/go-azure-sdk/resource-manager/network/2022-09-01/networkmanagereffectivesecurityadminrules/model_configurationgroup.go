package networkmanagereffectivesecurityadminrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationGroup struct {
	Id         *string                 `json:"id,omitempty"`
	Properties *NetworkGroupProperties `json:"properties,omitempty"`
}
