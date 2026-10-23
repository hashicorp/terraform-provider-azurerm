package indexers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndexingParametersConfiguration struct {
	AllowSkillsetToReadFileData                   *bool                                `json:"allowSkillsetToReadFileData,omitempty"`
	DataToExtract                                 *BlobIndexerDataToExtract            `json:"dataToExtract,omitempty"`
	DelimitedTextDelimiter                        *string                              `json:"delimitedTextDelimiter,omitempty"`
	DelimitedTextHeaders                          *string                              `json:"delimitedTextHeaders,omitempty"`
	DocumentRoot                                  *string                              `json:"documentRoot,omitempty"`
	ExcludedFileNameExtensions                    *string                              `json:"excludedFileNameExtensions,omitempty"`
	ExecutionEnvironment                          *IndexerExecutionEnvironment         `json:"executionEnvironment,omitempty"`
	FailOnUnprocessableDocument                   *bool                                `json:"failOnUnprocessableDocument,omitempty"`
	FailOnUnsupportedContentType                  *bool                                `json:"failOnUnsupportedContentType,omitempty"`
	FirstLineContainsHeaders                      *bool                                `json:"firstLineContainsHeaders,omitempty"`
	ImageAction                                   *BlobIndexerImageAction              `json:"imageAction,omitempty"`
	IndexStorageMetadataOnlyForOversizedDocuments *bool                                `json:"indexStorageMetadataOnlyForOversizedDocuments,omitempty"`
	IndexedFileNameExtensions                     *string                              `json:"indexedFileNameExtensions,omitempty"`
	ParsingMode                                   *BlobIndexerParsingMode              `json:"parsingMode,omitempty"`
	PdfTextRotationAlgorithm                      *BlobIndexerPDFTextRotationAlgorithm `json:"pdfTextRotationAlgorithm,omitempty"`
	QueryTimeout                                  *string                              `json:"queryTimeout,omitempty"`
}
