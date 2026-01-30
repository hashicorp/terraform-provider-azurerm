package certificate

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateFormat string

const (
	CertificateFormatCer CertificateFormat = "Cer"
	CertificateFormatPfx CertificateFormat = "Pfx"
)

func PossibleValuesForCertificateFormat() []string {
	return []string{
		string(CertificateFormatCer),
		string(CertificateFormatPfx),
	}
}

func (s *CertificateFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateFormat(input string) (*CertificateFormat, error) {
	vals := map[string]CertificateFormat{
		"cer": CertificateFormatCer,
		"pfx": CertificateFormatPfx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateFormat(input)
	return &out, nil
}

type CertificateProvisioningState string

const (
	CertificateProvisioningStateDeleting  CertificateProvisioningState = "Deleting"
	CertificateProvisioningStateFailed    CertificateProvisioningState = "Failed"
	CertificateProvisioningStateSucceeded CertificateProvisioningState = "Succeeded"
)

func PossibleValuesForCertificateProvisioningState() []string {
	return []string{
		string(CertificateProvisioningStateDeleting),
		string(CertificateProvisioningStateFailed),
		string(CertificateProvisioningStateSucceeded),
	}
}

func (s *CertificateProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateProvisioningState(input string) (*CertificateProvisioningState, error) {
	vals := map[string]CertificateProvisioningState{
		"deleting":  CertificateProvisioningStateDeleting,
		"failed":    CertificateProvisioningStateFailed,
		"succeeded": CertificateProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateProvisioningState(input)
	return &out, nil
}
