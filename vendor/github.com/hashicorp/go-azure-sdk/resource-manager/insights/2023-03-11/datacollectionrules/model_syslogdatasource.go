package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SyslogDataSource struct {
	FacilityNames *[]KnownSyslogDataSourceFacilityNames `json:"facilityNames,omitempty"`
	LogLevels     *[]KnownSyslogDataSourceLogLevels     `json:"logLevels,omitempty"`
	Name          *string                               `json:"name,omitempty"`
	Streams       *[]KnownSyslogDataSourceStreams       `json:"streams,omitempty"`
	TransformKql  *string                               `json:"transformKql,omitempty"`
}
