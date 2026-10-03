package indexers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobIndexerDataToExtract string

const (
	BlobIndexerDataToExtractAllMetadata        BlobIndexerDataToExtract = "allMetadata"
	BlobIndexerDataToExtractContentAndMetadata BlobIndexerDataToExtract = "contentAndMetadata"
	BlobIndexerDataToExtractStorageMetadata    BlobIndexerDataToExtract = "storageMetadata"
)

func PossibleValuesForBlobIndexerDataToExtract() []string {
	return []string{
		string(BlobIndexerDataToExtractAllMetadata),
		string(BlobIndexerDataToExtractContentAndMetadata),
		string(BlobIndexerDataToExtractStorageMetadata),
	}
}

func (s *BlobIndexerDataToExtract) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlobIndexerDataToExtract(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlobIndexerDataToExtract(input string) (*BlobIndexerDataToExtract, error) {
	vals := map[string]BlobIndexerDataToExtract{
		"allmetadata":        BlobIndexerDataToExtractAllMetadata,
		"contentandmetadata": BlobIndexerDataToExtractContentAndMetadata,
		"storagemetadata":    BlobIndexerDataToExtractStorageMetadata,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlobIndexerDataToExtract(input)
	return &out, nil
}

type BlobIndexerImageAction string

const (
	BlobIndexerImageActionGenerateNormalizedImagePerPage BlobIndexerImageAction = "generateNormalizedImagePerPage"
	BlobIndexerImageActionGenerateNormalizedImages       BlobIndexerImageAction = "generateNormalizedImages"
	BlobIndexerImageActionNone                           BlobIndexerImageAction = "none"
)

func PossibleValuesForBlobIndexerImageAction() []string {
	return []string{
		string(BlobIndexerImageActionGenerateNormalizedImagePerPage),
		string(BlobIndexerImageActionGenerateNormalizedImages),
		string(BlobIndexerImageActionNone),
	}
}

func (s *BlobIndexerImageAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlobIndexerImageAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlobIndexerImageAction(input string) (*BlobIndexerImageAction, error) {
	vals := map[string]BlobIndexerImageAction{
		"generatenormalizedimageperpage": BlobIndexerImageActionGenerateNormalizedImagePerPage,
		"generatenormalizedimages":       BlobIndexerImageActionGenerateNormalizedImages,
		"none":                           BlobIndexerImageActionNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlobIndexerImageAction(input)
	return &out, nil
}

type BlobIndexerPDFTextRotationAlgorithm string

const (
	BlobIndexerPDFTextRotationAlgorithmDetectAngles BlobIndexerPDFTextRotationAlgorithm = "detectAngles"
	BlobIndexerPDFTextRotationAlgorithmNone         BlobIndexerPDFTextRotationAlgorithm = "none"
)

func PossibleValuesForBlobIndexerPDFTextRotationAlgorithm() []string {
	return []string{
		string(BlobIndexerPDFTextRotationAlgorithmDetectAngles),
		string(BlobIndexerPDFTextRotationAlgorithmNone),
	}
}

func (s *BlobIndexerPDFTextRotationAlgorithm) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlobIndexerPDFTextRotationAlgorithm(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlobIndexerPDFTextRotationAlgorithm(input string) (*BlobIndexerPDFTextRotationAlgorithm, error) {
	vals := map[string]BlobIndexerPDFTextRotationAlgorithm{
		"detectangles": BlobIndexerPDFTextRotationAlgorithmDetectAngles,
		"none":         BlobIndexerPDFTextRotationAlgorithmNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlobIndexerPDFTextRotationAlgorithm(input)
	return &out, nil
}

type BlobIndexerParsingMode string

const (
	BlobIndexerParsingModeDefault       BlobIndexerParsingMode = "default"
	BlobIndexerParsingModeDelimitedText BlobIndexerParsingMode = "delimitedText"
	BlobIndexerParsingModeJson          BlobIndexerParsingMode = "json"
	BlobIndexerParsingModeJsonArray     BlobIndexerParsingMode = "jsonArray"
	BlobIndexerParsingModeJsonLines     BlobIndexerParsingMode = "jsonLines"
	BlobIndexerParsingModeText          BlobIndexerParsingMode = "text"
)

func PossibleValuesForBlobIndexerParsingMode() []string {
	return []string{
		string(BlobIndexerParsingModeDefault),
		string(BlobIndexerParsingModeDelimitedText),
		string(BlobIndexerParsingModeJson),
		string(BlobIndexerParsingModeJsonArray),
		string(BlobIndexerParsingModeJsonLines),
		string(BlobIndexerParsingModeText),
	}
}

func (s *BlobIndexerParsingMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlobIndexerParsingMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlobIndexerParsingMode(input string) (*BlobIndexerParsingMode, error) {
	vals := map[string]BlobIndexerParsingMode{
		"default":       BlobIndexerParsingModeDefault,
		"delimitedtext": BlobIndexerParsingModeDelimitedText,
		"json":          BlobIndexerParsingModeJson,
		"jsonarray":     BlobIndexerParsingModeJsonArray,
		"jsonlines":     BlobIndexerParsingModeJsonLines,
		"text":          BlobIndexerParsingModeText,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlobIndexerParsingMode(input)
	return &out, nil
}

type IndexerExecutionEnvironment string

const (
	IndexerExecutionEnvironmentPrivate  IndexerExecutionEnvironment = "private"
	IndexerExecutionEnvironmentStandard IndexerExecutionEnvironment = "standard"
)

func PossibleValuesForIndexerExecutionEnvironment() []string {
	return []string{
		string(IndexerExecutionEnvironmentPrivate),
		string(IndexerExecutionEnvironmentStandard),
	}
}

func (s *IndexerExecutionEnvironment) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIndexerExecutionEnvironment(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIndexerExecutionEnvironment(input string) (*IndexerExecutionEnvironment, error) {
	vals := map[string]IndexerExecutionEnvironment{
		"private":  IndexerExecutionEnvironmentPrivate,
		"standard": IndexerExecutionEnvironmentStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IndexerExecutionEnvironment(input)
	return &out, nil
}

type IndexerExecutionStatus string

const (
	IndexerExecutionStatusInProgress       IndexerExecutionStatus = "inProgress"
	IndexerExecutionStatusReset            IndexerExecutionStatus = "reset"
	IndexerExecutionStatusSuccess          IndexerExecutionStatus = "success"
	IndexerExecutionStatusTransientFailure IndexerExecutionStatus = "transientFailure"
)

func PossibleValuesForIndexerExecutionStatus() []string {
	return []string{
		string(IndexerExecutionStatusInProgress),
		string(IndexerExecutionStatusReset),
		string(IndexerExecutionStatusSuccess),
		string(IndexerExecutionStatusTransientFailure),
	}
}

func (s *IndexerExecutionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIndexerExecutionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIndexerExecutionStatus(input string) (*IndexerExecutionStatus, error) {
	vals := map[string]IndexerExecutionStatus{
		"inprogress":       IndexerExecutionStatusInProgress,
		"reset":            IndexerExecutionStatusReset,
		"success":          IndexerExecutionStatusSuccess,
		"transientfailure": IndexerExecutionStatusTransientFailure,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IndexerExecutionStatus(input)
	return &out, nil
}

type IndexerStatus string

const (
	IndexerStatusError   IndexerStatus = "error"
	IndexerStatusRunning IndexerStatus = "running"
	IndexerStatusUnknown IndexerStatus = "unknown"
)

func PossibleValuesForIndexerStatus() []string {
	return []string{
		string(IndexerStatusError),
		string(IndexerStatusRunning),
		string(IndexerStatusUnknown),
	}
}

func (s *IndexerStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIndexerStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIndexerStatus(input string) (*IndexerStatus, error) {
	vals := map[string]IndexerStatus{
		"error":   IndexerStatusError,
		"running": IndexerStatusRunning,
		"unknown": IndexerStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IndexerStatus(input)
	return &out, nil
}

type Prefer string

const (
	PreferReturnRepresentation Prefer = "return=representation"
)

func PossibleValuesForPrefer() []string {
	return []string{
		string(PreferReturnRepresentation),
	}
}

func (s *Prefer) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrefer(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrefer(input string) (*Prefer, error) {
	vals := map[string]Prefer{
		"return=representation": PreferReturnRepresentation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Prefer(input)
	return &out, nil
}
