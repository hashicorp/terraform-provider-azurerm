package list_upgrade

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const regPopulated = `package monitor

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

type Registration struct{}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		MonitorMetricAlertListResource{},
		MonitorScheduledQueryRulesAlertListResource{},
	}
}
`

const regEmpty = `package monitor

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

type Registration struct{}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
`

const regNoMethod = `package monitor

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

type Registration struct{}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}
`

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "registration.go")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp: %v", err)
	}
	return path
}

func TestRegisterListResource_InsertSorted(t *testing.T) {
	path := writeTemp(t, regPopulated)
	out, changed, err := RegisterListResource(path, "WorkspaceListResource")
	if err != nil {
		t.Fatalf("RegisterListResource: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed")
	}
	formatted, err := format.Source(out)
	if err != nil {
		t.Fatalf("not valid Go: %v\n%s", err, out)
	}
	got := string(formatted)
	mustContain(t, got, "WorkspaceListResource{},")

	// W sorts after both M entries, so it must be last in the slice.
	iMetric := strings.Index(got, "MonitorMetricAlertListResource{}")
	iWorkspace := strings.Index(got, "WorkspaceListResource{}")
	if iWorkspace < iMetric {
		t.Errorf("expected WorkspaceListResource to be ordered after the Monitor entries")
	}
}

func TestRegisterListResource_InsertFirst(t *testing.T) {
	path := writeTemp(t, regPopulated)
	out, _, err := RegisterListResource(path, "AaaListResource")
	if err != nil {
		t.Fatalf("RegisterListResource: %v", err)
	}
	got := formatGo(t, out)
	iAaa := strings.Index(got, "AaaListResource{}")
	iMetric := strings.Index(got, "MonitorMetricAlertListResource{}")
	if iAaa < 0 || iAaa > iMetric {
		t.Errorf("expected AaaListResource to be ordered first")
	}
}

func TestRegisterListResource_Empty(t *testing.T) {
	path := writeTemp(t, regEmpty)
	out, changed, err := RegisterListResource(path, "WorkspaceListResource")
	if err != nil {
		t.Fatalf("RegisterListResource: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed")
	}
	got := formatGo(t, out)
	mustContain(t, got, "return []sdk.FrameworkListWrappedResource{")
	mustContain(t, got, "WorkspaceListResource{},")
}

func TestRegisterListResource_CreateMethod(t *testing.T) {
	path := writeTemp(t, regNoMethod)
	out, changed, err := RegisterListResource(path, "WorkspaceListResource")
	if err != nil {
		t.Fatalf("RegisterListResource: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed")
	}
	got := formatGo(t, out)
	mustContain(t, got, "func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {")
	mustContain(t, got, "WorkspaceListResource{},")
}

func TestRegisterListResource_Idempotent(t *testing.T) {
	path := writeTemp(t, regPopulated)
	_, changed, err := RegisterListResource(path, "MonitorMetricAlertListResource")
	if err != nil {
		t.Fatalf("RegisterListResource: %v", err)
	}
	if changed {
		t.Errorf("expected no change when the list resource is already registered")
	}
}

func formatGo(t *testing.T, src []byte) string {
	t.Helper()
	out, err := format.Source(src)
	if err != nil {
		t.Fatalf("not valid Go: %v\n%s", err, src)
	}
	return string(out)
}
