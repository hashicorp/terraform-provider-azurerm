package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateIssuerItemOperationPredicate struct {
	Id       *string
	Provider *string
}

func (p CertificateIssuerItemOperationPredicate) Matches(input CertificateIssuerItem) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Provider != nil && (input.Provider == nil || *p.Provider != *input.Provider) {
		return false
	}

	return true
}

type CertificateItemOperationPredicate struct {
	Id  *string
	X5t *string
}

func (p CertificateItemOperationPredicate) Matches(input CertificateItem) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.X5t != nil && (input.X5t == nil || *p.X5t != *input.X5t) {
		return false
	}

	return true
}
