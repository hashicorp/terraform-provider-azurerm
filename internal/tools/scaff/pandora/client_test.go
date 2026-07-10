package pandora

import (
	"net/http"
	"testing"
	"time"
)

// serverAvailable returns true if the Pandora Data API is reachable, so live
// tests can skip cleanly when it isn't running.
func serverAvailable(baseURL string) bool {
	c := &http.Client{Timeout: 2 * time.Second}
	resp, err := c.Get(baseURL + resourceManagerBasePath)
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// TestLive_RedHatOpenShift exercises the full client against a running Pandora
// Data API using the RedHatOpenShift service as a fixture.
func TestLive_RedHatOpenShift(t *testing.T) {
	client := NewClient(DefaultBaseURL)
	if !serverAvailable(client.BaseURL) {
		t.Skipf("pandora data api not reachable at %s; skipping live test", client.BaseURL)
	}

	services, err := client.ListServices()
	if err != nil {
		t.Fatalf("ListServices: %v", err)
	}
	if _, ok := services["RedHatOpenShift"]; !ok {
		t.Fatalf("expected RedHatOpenShift in services list")
	}

	svc, err := client.GetService("RedHatOpenShift")
	if err != nil {
		t.Fatalf("GetService: %v", err)
	}
	if svc.ResourceProvider != "Microsoft.RedHatOpenShift" {
		t.Fatalf("unexpected resourceProvider %q", svc.ResourceProvider)
	}
	if _, ok := svc.Versions["2025-07-25"]; !ok {
		t.Fatalf("expected version 2025-07-25")
	}

	ver, err := client.GetVersion("RedHatOpenShift", "2025-07-25")
	if err != nil {
		t.Fatalf("GetVersion: %v", err)
	}
	res, ok := ver.Resources["OpenShiftClusters"]
	if !ok {
		t.Fatalf("expected OpenShiftClusters resource")
	}
	if res.SchemaURI == "" || res.OperationsURI == "" {
		t.Fatalf("expected schema and operations URIs to be populated")
	}

	schema, err := client.GetResourceSchema("RedHatOpenShift", "2025-07-25", "OpenShiftClusters")
	if err != nil {
		t.Fatalf("GetResourceSchema: %v", err)
	}
	if _, ok := schema.Models["OpenShiftCluster"]; !ok {
		t.Fatalf("expected OpenShiftCluster model")
	}
	if _, ok := schema.Constants["Visibility"]; !ok {
		t.Fatalf("expected Visibility constant")
	}
	id, ok := schema.ResourceIds["OpenShiftClusterId"]
	if !ok {
		t.Fatalf("expected OpenShiftClusterId resource id")
	}
	if len(id.Segments) == 0 {
		t.Fatalf("expected id segments")
	}

	ops, err := client.GetResourceOperations("RedHatOpenShift", "2025-07-25", "OpenShiftClusters")
	if err != nil {
		t.Fatalf("GetResourceOperations: %v", err)
	}
	create, ok := ops.Operations["CreateOrUpdate"]
	if !ok {
		t.Fatalf("expected CreateOrUpdate operation")
	}
	if create.Method != "PUT" || !create.LongRunning {
		t.Fatalf("expected CreateOrUpdate to be a long-running PUT, got method=%s lro=%v", create.Method, create.LongRunning)
	}
	if create.RequestObject == nil || create.RequestObject.ReferenceName == nil || *create.RequestObject.ReferenceName != "OpenShiftCluster" {
		t.Fatalf("expected CreateOrUpdate request object OpenShiftCluster")
	}
}
