// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package features

type Provider interface {
	UserFeatures() UserFeatures
}

func From(meta any) UserFeatures {
	return meta.(Provider).UserFeatures()
}
