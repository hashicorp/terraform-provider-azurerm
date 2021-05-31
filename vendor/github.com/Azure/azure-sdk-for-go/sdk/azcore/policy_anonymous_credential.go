// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

func anonCredAuthPolicyFunc(AuthenticationPolicyOptions) Policy {
	return PolicyFunc(anonCredPolicyFunc)
}

func anonCredPolicyFunc(req *Request) (*Response, error) {
	return req.Next()
}

// AnonymousCredential is for use with HTTP(S) requests that read public resource
// or for use with Shared Access Signatures (SAS).
func AnonymousCredential() Credential {
	return credentialFunc(anonCredAuthPolicyFunc)
}
