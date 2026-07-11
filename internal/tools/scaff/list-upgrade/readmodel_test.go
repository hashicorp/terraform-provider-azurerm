package list_upgrade

import (
	"os"
	"path/filepath"
	"testing"
)

func TestModelFieldType(t *testing.T) {
	src := `type VirtualHubsGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *VirtualHub
}

func somethingElse() {}
`
	if got := modelFieldType(src); got != "VirtualHub" {
		t.Errorf("modelFieldType = %q, want VirtualHub", got)
	}
}

func TestReadModelFromGetMethod(t *testing.T) {
	dir := t.TempDir()
	method := `package virtualwans

type VirtualHubIPConfigurationGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *HubIPConfiguration
}
`
	if err := os.WriteFile(filepath.Join(dir, "method_virtualhubipconfigurationget.go"), []byte(method), 0o644); err != nil {
		t.Fatalf("writing method file: %v", err)
	}
	if got := readModelFromGetMethod(dir, "VirtualHubIPConfigurationGet"); got != "HubIPConfiguration" {
		t.Errorf("readModelFromGetMethod = %q, want HubIPConfiguration", got)
	}
	if got := readModelFromGetMethod(dir, "MissingGet"); got != "" {
		t.Errorf("expected empty for missing method, got %q", got)
	}
}
