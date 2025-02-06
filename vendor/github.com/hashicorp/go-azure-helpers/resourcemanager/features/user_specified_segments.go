// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package features

// TreatUserSpecifiedSegmentsAsCaseInsensitive is a feature-toggle which specifies whether User Specified
// Resource ID Segments should be compared case-insensitively as required.
//
// @tombuildsstuff: whilst this IS EXPOSED in the public interface - this is NOT READY FOR USE and should
// not be exposed to users (i.e. as a feature-toggle) until this work is completed - as this'll become a source of knock-on problems
// rather than being useful.
//
// There are a number of dependencies to enabling this, including completing the standardiation on the
// `ResourceId` interface and the `ResourceIDReference` schema types - and surrounding updates.
var TreatUserSpecifiedSegmentsAsCaseInsensitive = false
