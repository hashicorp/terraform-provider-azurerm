package list_upgrade

import (
	"go/format"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnalyze_UntypedWithIdentity(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "untyped_with_identity.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	assertEqual(t, "kind", r.Kind.String(), "untyped")
	assertEqual(t, "constructor", r.ConstructorFunc, "resourceVirtualHubIP")
	assertEqual(t, "base name", r.BaseName, "VirtualHubIP")
	assertEqual(t, "read func", r.ReadFunc, "resourceVirtualHubIPRead")
	assertEqual(t, "create func", r.CreateFunc, "resourceVirtualHubIPCreate")
	assertEqual(t, "service", r.ServiceName, "Network")
	assertEqual(t, "client field", r.ClientField, "VirtualWANs")
	assertEqual(t, "id package", r.IDPackage, "commonids")
	assertEqual(t, "id type", r.IDTypeName, "VirtualHubIPConfigurationId")
	assertEqual(t, "id parse", r.IDParseFunc, "ParseVirtualHubIPConfigurationID")
	assertEqual(t, "sdk package", r.SDKPackage, "virtualwans")
	assertEqual(t, "get method", r.GetMethod, "VirtualHubIPConfigurationGet")

	if !r.HasIdentity {
		t.Errorf("expected HasIdentity to be true")
	}
	if r.HasFlatten {
		t.Errorf("expected HasFlatten to be false")
	}
}

func TestAnalyze_UntypedChildConnection_ParentDetection(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "untyped_child_connection.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	// The child ID base is not a prefix of the parent, and a second unrelated ID
	// (remote_virtual_network_id) is parsed in the same func, so detection must
	// rely on which parent ID constructs the child ID.
	assertEqual(t, "id base", r.IDBase, "HubVirtualNetworkConnection")
	assertEqual(t, "parent id type", r.ParentIDType, "VirtualHubId")
	assertEqual(t, "parent attr", r.ParentAttr, "virtual_hub_id")
	assertEqual(t, "parent parse", r.ParentParseFunc, "ParseVirtualHubID")
	assertEqual(t, "list method", r.ListMethod, "HubVirtualNetworkConnectionsList")
}

func TestAnalyze_TypedChildRoutingIntent_ParentDetection(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "typed_child_routing_intent.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	// A typed child parses its parent from a decoded model field
	// (model.VirtualHubId) inside the Create closure; detection must still find
	// the parent that constructs the child ID and resolve the tfschema attribute.
	assertEqual(t, "kind", r.Kind.String(), "typed")
	assertEqual(t, "id base", r.IDBase, "RoutingIntent")
	assertEqual(t, "parent id type", r.ParentIDType, "VirtualHubId")
	assertEqual(t, "parent attr", r.ParentAttr, "virtual_hub_id")
	assertEqual(t, "parent parse", r.ParentParseFunc, "ParseVirtualHubID")
	assertEqual(t, "get method", r.GetMethod, "RoutingIntentGet")
	assertEqual(t, "list method", r.ListMethod, "RoutingIntentList")
}

func TestAnalyze_UntypedSubnet_SDKDerivedParent(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "untyped_subnet.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	// subnet imports three go-azure-sdk packages; the Get call's subnets-qualified
	// options argument must disambiguate the model+client package.
	assertEqual(t, "sdk package", r.SDKPackage, "subnets")

	// The parent scope is not visible in the source (the SubnetId is built from
	// config strings), so it must be derived from the vendored SDK list method
	// subnets.ListComplete(ctx, commonids.VirtualNetworkId).
	assertEqual(t, "parent id type", r.ParentIDType, "VirtualNetworkId")
	assertEqual(t, "parent package", r.ParentPackage, "commonids")
	assertEqual(t, "parent attr", r.ParentAttr, "virtual_network_id")
	assertEqual(t, "parent parse", r.ParentParseFunc, "ParseVirtualNetworkID")
	assertEqual(t, "parent validate", r.ParentValidateFunc, "ValidateVirtualNetworkID")
	assertEqual(t, "list method", r.ListMethod, "List")
	assertEqual(t, "read model", r.ReadModel, "Subnet")

	// The pre-existing flatten takes its id by value, so the list generator must
	// dereference the parsed (pointer) id.
	if !r.FlattenIDValue {
		t.Errorf("expected FlattenIDValue to be true for a value-id flatten")
	}
}

func TestUpgrade_UntypedExtractsFlatten(t *testing.T) {
	r, err := Analyze(filepath.Join("testdata", "untyped_with_identity.go"))
	if err != nil {
		t.Fatalf("Analyze: %v", err)
	}

	out, changed, err := r.Upgrade(UpgradeOptions{
		ExtractFlatten: true,
		ReadModel:      "VirtualHubIPConfiguration",
		ResourceName:   "virtual_hub_ip",
	})
	if err != nil {
		t.Fatalf("Upgrade: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed to be true")
	}

	formatted, err := format.Source(out)
	if err != nil {
		t.Fatalf("result is not valid Go: %v\n\n%s", err, out)
	}
	got := string(formatted)

	mustContain(t, got, "func resourceVirtualHubIPFlatten(d *pluginsdk.ResourceData, id *commonids.VirtualHubIPConfigurationId, model *virtualwans.VirtualHubIPConfiguration) error {")
	mustContain(t, got, "return resourceVirtualHubIPFlatten(d, id, resp.Model)")
	mustContain(t, got, "return pluginsdk.SetResourceIdentityData(d, id)")
	mustContain(t, got, "if model != nil {")
	if strings.Contains(got, "if model := resp.Model; model != nil") {
		t.Errorf("expected the resp.Model guard to be collapsed in flatten")
	}
}
