package dpscertificate

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatePurpose string

const (
	CertificatePurposeClientAuthentication CertificatePurpose = "clientAuthentication"
	CertificatePurposeServerAuthentication CertificatePurpose = "serverAuthentication"
)

func PossibleValuesForCertificatePurpose() []string {
	return []string{
		string(CertificatePurposeClientAuthentication),
		string(CertificatePurposeServerAuthentication),
	}
}

func parseCertificatePurpose(input string) (*CertificatePurpose, error) {
	vals := map[string]CertificatePurpose{
		"clientauthentication": CertificatePurposeClientAuthentication,
		"serverauthentication": CertificatePurposeServerAuthentication,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificatePurpose(input)
	return &out, nil
}
