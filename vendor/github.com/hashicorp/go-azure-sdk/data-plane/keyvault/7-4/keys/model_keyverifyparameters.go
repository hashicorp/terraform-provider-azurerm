package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVerifyParameters struct {
	Alg    JsonWebKeySignatureAlgorithm `json:"alg"`
	Digest string                       `json:"digest"`
	Value  string                       `json:"value"`
}
