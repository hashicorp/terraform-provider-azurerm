package networkvirtualappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkVirtualApplianceBootDiagnosticParameters struct {
	ConsoleScreenshotStorageSasURL *string `json:"consoleScreenshotStorageSasUrl,omitempty"`
	InstanceId                     *int64  `json:"instanceId,omitempty"`
	SerialConsoleStorageSasURL     *string `json:"serialConsoleStorageSasUrl,omitempty"`
}
