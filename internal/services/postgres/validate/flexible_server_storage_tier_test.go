// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/servers"
)

func TestCompareAzureManagedDiskPerformance(t *testing.T) {
	p10 := servers.AzureManagedDiskPerformanceTiers("p10")
	p20 := servers.AzureManagedDiskPerformanceTiers("p20")

	notFaster := CompareAzureManagedDiskPerformance(p10, p20)
	if notFaster {
		t.Errorf("expected no error, %s is slower than %s - got: %v", p10, p20, notFaster)
	}

	notFaster = CompareAzureManagedDiskPerformance(p20, p20)
	if notFaster {
		t.Errorf("expected no error, %s is slower/equal than %s - got: %v", p20, p20, notFaster)
	}

	isFaster := CompareAzureManagedDiskPerformance(p20, p10)
	if !isFaster {
		t.Errorf("expected no error, %s is slower than %s - got: %v", p20, p10, isFaster)
	}
}
