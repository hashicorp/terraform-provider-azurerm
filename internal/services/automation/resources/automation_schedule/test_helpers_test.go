// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_schedule_test

import (
	"os"
	"path/filepath"
)

var (
	testCertThumbprintRaw, _ = os.ReadFile(filepath.Join("testdata", "automation_certificate_test.thumb"))
	testCertRaw, _           = os.ReadFile(filepath.Join("testdata", "automation_certificate_test.pfx"))
)
