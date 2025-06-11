package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WindowsEventLogDataSource struct {
	Name         *string                                  `json:"name,omitempty"`
	Streams      *[]KnownWindowsEventLogDataSourceStreams `json:"streams,omitempty"`
	TransformKql *string                                  `json:"transformKql,omitempty"`
	XPathQueries *[]string                                `json:"xPathQueries,omitempty"`
}
