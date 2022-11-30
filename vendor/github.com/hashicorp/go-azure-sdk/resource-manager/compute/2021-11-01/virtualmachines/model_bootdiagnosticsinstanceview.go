package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BootDiagnosticsInstanceView struct {
	ConsoleScreenshotBlobUri *string             `json:"consoleScreenshotBlobUri,omitempty"`
	SerialConsoleLogBlobUri  *string             `json:"serialConsoleLogBlobUri,omitempty"`
	Status                   *InstanceViewStatus `json:"status,omitempty"`
}
