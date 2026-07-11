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

func TestParseListCompleteMethod(t *testing.T) {
	cases := []struct {
		name     string
		line     string
		prefix   string
		wantName string
		wantID   string
		wantPkg  string
	}{
		{
			name:     "subscription scoped",
			line:     "func (c VirtualWANsClient) VirtualHubsListComplete(ctx context.Context, id commonids.SubscriptionId) (result VirtualHubsListCompleteResult, err error) {",
			prefix:   "VirtualHubs",
			wantName: "VirtualHubsList",
			wantID:   "SubscriptionId",
			wantPkg:  "commonids",
		},
		{
			name:     "resource group scoped",
			line:     "func (c VirtualWANsClient) VirtualHubsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (result VirtualHubsListByResourceGroupCompleteResult, err error) {",
			prefix:   "VirtualHubs",
			wantName: "VirtualHubsListByResourceGroup",
			wantID:   "ResourceGroupId",
			wantPkg:  "commonids",
		},
		{
			name:     "parent scoped",
			line:     "func (c SubnetsClient) ListComplete(ctx context.Context, id commonids.VirtualNetworkId) (ListCompleteResult, error) {",
			prefix:   "",
			wantName: "List",
			wantID:   "VirtualNetworkId",
			wantPkg:  "commonids",
		},
		{
			name:   "matching-predicate variant ignored",
			line:   "func (c SubnetsClient) ListCompleteMatchingPredicate(ctx context.Context, id commonids.VirtualNetworkId, predicate SubnetOperationPredicate) (result ListCompleteResult, err error) {",
			prefix: "",
		},
		{
			name:   "wrong prefix",
			line:   "func (c VirtualWANsClient) VpnGatewaysListComplete(ctx context.Context, id commonids.SubscriptionId) (result VpnGatewaysListCompleteResult, err error) {",
			prefix: "VirtualHubs",
		},
		{
			name:   "not a list method",
			line:   "func (c VirtualWANsClient) VirtualHubsGetComplete(ctx context.Context, id commonids.SubscriptionId) (result VirtualHubsGetCompleteResult, err error) {",
			prefix: "VirtualHubs",
		},
		{
			name:   "get is not a list method",
			line:   "func (c VirtualWANsClient) VirtualHubsGet(ctx context.Context, id commonids.VirtualHubId) (result VirtualHubsGetOperationResponse, err error) {",
			prefix: "VirtualHubs",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotName, gotID, gotPkg := parseListCompleteMethod(tc.line, tc.prefix)
			if gotName != tc.wantName || gotID != tc.wantID || gotPkg != tc.wantPkg {
				t.Errorf("parseListCompleteMethod = (%q, %q, %q), want (%q, %q, %q)", gotName, gotID, gotPkg, tc.wantName, tc.wantID, tc.wantPkg)
			}
		})
	}
}

func TestListMethodsFromVendor(t *testing.T) {
	dir := t.TempDir()
	src := `package virtualwans

func (c VirtualWANsClient) VirtualHubsListComplete(ctx context.Context, id commonids.SubscriptionId) (result VirtualHubsListCompleteResult, err error) {
	return
}

func (c VirtualWANsClient) VirtualHubsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (result VirtualHubsListByResourceGroupCompleteResult, err error) {
	return
}
`
	if err := os.WriteFile(filepath.Join(dir, "client.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("writing client file: %v", err)
	}
	m := listMethodsFromVendor(dir, "VirtualHubs")
	if m.sub != "VirtualHubsList" {
		t.Errorf("subscription method = %q, want VirtualHubsList", m.sub)
	}
	if m.rg != "VirtualHubsListByResourceGroup" {
		t.Errorf("resource-group method = %q, want VirtualHubsListByResourceGroup", m.rg)
	}
}

func TestListMethodsFromVendor_ParentScoped(t *testing.T) {
	dir := t.TempDir()
	src := `package subnets

func (c SubnetsClient) List(ctx context.Context, id commonids.VirtualNetworkId) (result ListOperationResponse, err error) {
	return
}

func (c SubnetsClient) ListComplete(ctx context.Context, id commonids.VirtualNetworkId) (ListCompleteResult, error) {
	return
}

func (c SubnetsClient) ListCompleteMatchingPredicate(ctx context.Context, id commonids.VirtualNetworkId, predicate SubnetOperationPredicate) (result ListCompleteResult, err error) {
	return
}
`
	if err := os.WriteFile(filepath.Join(dir, "method_list.go"), []byte(src), 0o644); err != nil {
		t.Fatalf("writing method file: %v", err)
	}
	m := listMethodsFromVendor(dir, "")
	if m.sub != "" || m.rg != "" {
		t.Errorf("expected no subscription/resource-group methods, got sub=%q rg=%q", m.sub, m.rg)
	}
	if m.parent != "List" {
		t.Errorf("parent method = %q, want List", m.parent)
	}
	if m.parentIDType != "VirtualNetworkId" {
		t.Errorf("parent id type = %q, want VirtualNetworkId", m.parentIDType)
	}
	if m.parentIDPkg != "commonids" {
		t.Errorf("parent id pkg = %q, want commonids", m.parentIDPkg)
	}
}
