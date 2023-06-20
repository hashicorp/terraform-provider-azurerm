package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSettingsProperties struct {
	NetworkAdapters *[]NetworkAdapter `json:"networkAdapters,omitempty"`
}
