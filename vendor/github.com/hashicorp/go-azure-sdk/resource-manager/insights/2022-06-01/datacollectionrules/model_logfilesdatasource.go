package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogFilesDataSource struct {
	FilePatterns []string                      `json:"filePatterns"`
	Format       KnownLogFilesDataSourceFormat `json:"format"`
	Name         *string                       `json:"name,omitempty"`
	Settings     *LogFileSettings              `json:"settings,omitempty"`
	Streams      []string                      `json:"streams"`
}
