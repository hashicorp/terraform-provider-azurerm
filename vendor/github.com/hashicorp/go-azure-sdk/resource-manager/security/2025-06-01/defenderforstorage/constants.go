package defenderforstorage

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobScanResultsOptions string

const (
	BlobScanResultsOptionsBlobIndexTags BlobScanResultsOptions = "blobIndexTags"
	BlobScanResultsOptionsNone          BlobScanResultsOptions = "None"
)

func PossibleValuesForBlobScanResultsOptions() []string {
	return []string{
		string(BlobScanResultsOptionsBlobIndexTags),
		string(BlobScanResultsOptionsNone),
	}
}

func (s *BlobScanResultsOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBlobScanResultsOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBlobScanResultsOptions(input string) (*BlobScanResultsOptions, error) {
	vals := map[string]BlobScanResultsOptions{
		"blobindextags": BlobScanResultsOptionsBlobIndexTags,
		"none":          BlobScanResultsOptionsNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BlobScanResultsOptions(input)
	return &out, nil
}
