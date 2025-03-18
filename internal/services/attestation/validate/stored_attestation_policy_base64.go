// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackofallops/kermit/sdk/attestation/2022-08-01/attestation"
)

func ContainsABase64UriEncodedJWTOfAStoredAttestationPolicy(value interface{}, key string) (warnings []string, errs []error) {
	v, ok := value.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("%q cannot be an empty string", key))
		return
	}

	split := strings.Split(v, ".")
	if len(split) != 3 {
		errs = append(errs, fmt.Errorf("expected %q to be a JWT with 3 segments but got %d segments", key, len(split)))
		return
	}

	// decode the JWT into a StoredAttestationPolicy object
	decodedJwtSegment, err := base64.RawURLEncoding.DecodeString(split[1])
	if err != nil {
		errs = append(errs, fmt.Errorf("base64-decoding the first JWT Segment %q: %+v", split[1], err))
		return
	}
	var result attestation.StoredAttestationPolicy
	if err := json.Unmarshal(decodedJwtSegment, &result); err != nil {
		errs = append(errs, fmt.Errorf("unmarshaling into StoredAttestationPolicy - please check your policy against the documentation at https://learn.microsoft.com/azure/attestation/author-sign-policy: %+v", err))
		return
	}
	if result.AttestationPolicy == nil {
		errs = append(errs, errors.New("expected a key for AttestationPolicy but didn't get one - see https://learn.microsoft.com/azure/attestation/author-sign-policy for more information"))
		return
	}

	return
}
