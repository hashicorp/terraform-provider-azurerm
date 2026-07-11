// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ir

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/pandora"
)

// Options are the inputs required to resolve a ResourceIR. Only a handful are
// required; everything else is derived deterministically from Pandora.
type Options struct {
	ARMType        string // e.g. "Microsoft.RedHatOpenShift/openShiftClusters"
	Service        string // optional explicit Pandora service name (overrides ARMType-derived)
	Resource       string // optional explicit Pandora resource key (overrides ARMType-derived)
	APIVersion     string // optional; defaults to the latest non-preview version
	Name           string // terraform resource name (snake), e.g. "redhat_openshift_cluster"
	GoName         string // optional Go/Camel identifier override, e.g. "RedHatOpenShiftCluster"; derived from Name when empty
	ServicePackage string // service package directory, e.g. "redhatopenshift"

	ProviderName      string // e.g. "azurerm"
	ProviderGithubOrg string // e.g. "hashicorp"
}

// Resolve builds a ResourceIR from the Pandora Data API for the given options.
func Resolve(client *pandora.Client, opts Options) (*ResourceIR, error) {
	service, version, resourceKey, svc, err := resolveTarget(client, opts)
	if err != nil {
		return nil, err
	}

	schema, err := client.GetResourceSchema(service, version, resourceKey)
	if err != nil {
		return nil, fmt.Errorf("fetching schema: %w", err)
	}
	ops, err := client.GetResourceOperations(service, version, resourceKey)
	if err != nil {
		return nil, fmt.Errorf("fetching operations: %w", err)
	}

	providerName := opts.ProviderName
	if providerName == "" {
		providerName = "azurerm"
	}
	org := opts.ProviderGithubOrg
	if org == "" {
		org = "hashicorp"
	}
	servicePackage := opts.ServicePackage
	if servicePackage == "" {
		servicePackage = strings.ToLower(service)
	}

	sdkPackage := strings.ToLower(resourceKey)
	name := opts.GoName
	if name == "" {
		name = camel(opts.Name)
	}

	res := &ResourceIR{
		ProviderName:      providerName,
		ProviderGithubOrg: org,
		ServicePackage:    servicePackage,
		ServiceName:       service,
		ResourceProvider:  svc.ResourceProvider,
		APIVersion:        version,
		PandoraResource:   resourceKey,

		Name:          name,
		TerraformType: fmt.Sprintf("%s_%s", providerName, snake(opts.Name)),

		SDKPackage:    sdkPackage,
		SDKImportPath: fmt.Sprintf("github.com/%s/go-azure-sdk/resource-manager/%s/%s/%s", org, strings.ToLower(service), version, sdkPackage),
		ClientField:   resourceKey + "Client",

		ModelStructName: name + "Model",
	}

	primaryID := classifyOperations(res, ops.Operations)
	if primaryID == "" {
		return nil, fmt.Errorf("unable to determine primary resource ID from operations")
	}
	classifyListOperations(res, ops.Operations)
	if err := deriveID(res, schema, primaryID); err != nil {
		return nil, err
	}

	if err := res.walkSchema(schema); err != nil {
		return nil, err
	}

	res.resolveUpdateModel(schema)

	return res, nil
}

// resolveTarget determines the Pandora service, version and resource key for the
// given options, following the ARM type where an explicit value is not provided.
func resolveTarget(client *pandora.Client, opts Options) (service, version, resourceKey string, svc *pandora.ServiceResponse, err error) {
	namespace := ""
	resPart := opts.Resource
	if opts.ARMType != "" {
		parts := strings.SplitN(opts.ARMType, "/", 2)
		namespace = parts[0]
		if len(parts) == 2 && resPart == "" {
			resPart = parts[1]
		}
	}

	service = opts.Service
	if service == "" {
		service = strings.TrimPrefix(namespace, "Microsoft.")
	}
	if service == "" {
		return "", "", "", nil, fmt.Errorf("unable to determine service: provide -arm-type or -service")
	}
	if resPart == "" {
		return "", "", "", nil, fmt.Errorf("unable to determine resource: provide -arm-type or -resource")
	}

	svc, err = client.GetService(service)
	if err != nil || (namespace != "" && svc != nil && svc.ResourceProvider != "" && !strings.EqualFold(svc.ResourceProvider, namespace)) {
		// Fall back to scanning services and matching on the resource provider.
		foundSvc, foundName, ferr := findServiceByRP(client, namespace)
		if ferr != nil {
			if err != nil {
				return "", "", "", nil, fmt.Errorf("resolving service %q: %w", service, err)
			}
			return "", "", "", nil, ferr
		}
		svc = foundSvc
		service = foundName
	}

	version = opts.APIVersion
	if version == "" {
		version = latestVersion(svc.Versions)
	}
	if version == "" {
		return "", "", "", nil, fmt.Errorf("no API versions available for service %q", service)
	}

	verResp, err := client.GetVersion(service, version)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("fetching version %q: %w", version, err)
	}
	// For nested (child) ARM types such as "storageMovers/agents", the Pandora
	// resource key is the final path segment ("Agents").
	matchName := resPart
	if idx := strings.LastIndex(resPart, "/"); idx >= 0 {
		matchName = resPart[idx+1:]
	}
	resourceKey = matchResourceKey(verResp.Resources, matchName)
	if resourceKey == "" {
		return "", "", "", nil, fmt.Errorf("resource %q not found in %s/%s", resPart, service, version)
	}

	return service, version, resourceKey, svc, nil
}

// matchResourceKey finds the Pandora resource key corresponding to an ARM type
// segment. It prefers an exact (case-insensitive) match, then tolerates
// singular/plural differences between the ARM segment and the Pandora key (e.g.
// ARM "automationAccounts" vs key "AutomationAccount"). Iteration is sorted so
// the result is deterministic.
func matchResourceKey(resources map[string]pandora.ResourceSummary, matchName string) string {
	keys := make([]string, 0, len(resources))
	for k := range resources {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if strings.EqualFold(k, matchName) {
			return k
		}
	}
	target := singularize(strings.ToLower(matchName))
	for _, k := range keys {
		if singularize(strings.ToLower(k)) == target {
			return k
		}
	}
	return ""
}

// findServiceByRP scans all services and returns the one whose resource provider
// namespace matches. Used as a fallback when the service name cannot be derived
// directly from the ARM type.
func findServiceByRP(client *pandora.Client, namespace string) (*pandora.ServiceResponse, string, error) {
	if namespace == "" {
		return nil, "", fmt.Errorf("no resource provider namespace to match against")
	}
	services, err := client.ListServices()
	if err != nil {
		return nil, "", err
	}
	names := make([]string, 0, len(services))
	for name := range services {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		svc, err := client.GetService(name)
		if err != nil {
			continue
		}
		if strings.EqualFold(svc.ResourceProvider, namespace) {
			return svc, name, nil
		}
	}
	return nil, "", fmt.Errorf("no service found for resource provider %q", namespace)
}

// latestVersion returns the most recent non-preview version, falling back to the
// most recent version if all are preview.
func latestVersion(versions map[string]pandora.VersionSummary) string {
	keys := make([]string, 0, len(versions))
	for k := range versions {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i := len(keys) - 1; i >= 0; i-- {
		if !versions[keys[i]].Preview {
			return keys[i]
		}
	}
	if len(keys) > 0 {
		return keys[len(keys)-1]
	}
	return ""
}

// classifyOperations maps the Pandora operations to CRUD roles by their
// conventional names, records LRO flags and request/response models, and returns
// the primary (single-instance) resource ID name.
func classifyOperations(res *ResourceIR, ops map[string]pandora.Operation) string {
	names := sortedOpNames(ops)

	if op, name := pickOp(ops, names, "PUT", []string{"CreateOrUpdate", "Create"}, nil); name != "" {
		res.CreateOp = name
		res.CreateLRO = op.LongRunning
		res.CreateModel = refName(op.RequestObject)
	}
	if op, name := pickOp(ops, names, "GET", []string{"Get"}, nil); name != "" {
		res.ReadOp = name
		res.ReadModel = refName(op.ResponseObject)
	}
	if op, name := pickOp(ops, names, "PATCH", []string{"Update"}, []string{"CreateOrUpdate"}); name != "" {
		res.Updatable = true
		res.UpdateOp = name
		res.UpdateLRO = op.LongRunning
		res.UpdateModel = refName(op.RequestObject)
	}
	if op, name := pickOp(ops, names, "DELETE", []string{"Delete"}, nil); name != "" {
		res.DeleteOp = name
		res.DeleteLRO = op.LongRunning
	}

	// Primary ID is the one used by the single-instance Read (Get); fall back to Create.
	if res.ReadOp != "" {
		if op, ok := ops[res.ReadOp]; ok && op.ResourceIDName != nil {
			return *op.ResourceIDName
		}
	}
	if res.CreateOp != "" {
		if op, ok := ops[res.CreateOp]; ok && op.ResourceIDName != nil {
			return *op.ResourceIDName
		}
	}
	return ""
}

// sortedOpNames returns the operation names sorted for deterministic iteration.
func sortedOpNames(ops map[string]pandora.Operation) []string {
	names := make([]string, 0, len(ops))
	for name := range ops {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// pickOp selects the operation implementing a CRUD verb. It first tries an exact
// name match (verbs are tried in order); failing that it returns the operation
// that uses the expected HTTP method and whose name ends with one of the verb
// suffixes - Azure commonly prefixes operations with the collection name (e.g.
// "AccountsGet", "AccountsCreateOrUpdate"). Operations whose name ends with any
// suffix in exclude are skipped so, for example, "Update" does not select a
// "CreateOrUpdate" operation. Iteration is deterministic (names pre-sorted).
func pickOp(ops map[string]pandora.Operation, names []string, method string, verbs, exclude []string) (pandora.Operation, string) {
	for _, v := range verbs {
		if op, ok := ops[v]; ok {
			return op, v
		}
	}
	for _, name := range names {
		op := ops[name]
		if !strings.EqualFold(op.Method, method) {
			continue
		}
		if hasSuffixFold(name, exclude) {
			continue
		}
		if hasSuffixFold(name, verbs) {
			return op, name
		}
	}
	return pandora.Operation{}, ""
}

// hasSuffixFold reports whether s ends with any of the given suffixes, ignoring case.
func hasSuffixFold(s string, suffixes []string) bool {
	lower := strings.ToLower(s)
	for _, suf := range suffixes {
		if suf != "" && strings.HasSuffix(lower, strings.ToLower(suf)) {
			return true
		}
	}
	return false
}

// classifyListOperations identifies the subscription- and resource-group-scoped
// list operations (used to generate a list resource) by their resource ID name,
// preferring the canonical operation names when several candidates exist.
func classifyListOperations(res *ResourceIR, ops map[string]pandora.Operation) {
	for _, name := range sortedOpNames(ops) {
		op := ops[name]
		if op.Method != "GET" || op.ResourceIDName == nil {
			continue
		}
		switch *op.ResourceIDName {
		case "SubscriptionId":
			// A subscription-scoped GET is unambiguously a list, regardless of the
			// operation name (which may be prefixed, e.g. "AccountsListBySubscription").
			if res.ListBySubscriptionOp == "" || name == "List" || strings.HasSuffix(name, "ListBySubscription") {
				res.ListBySubscriptionOp = name
			}
		case "ResourceGroupId":
			if res.ListByResourceGroupOp == "" || name == "ListByResourceGroup" || strings.HasSuffix(name, "ListByResourceGroup") {
				res.ListByResourceGroupOp = name
			}
		default:
			// A list scoped to something other than the subscription or resource
			// group is a parent-scoped (child resource) list. Guard on the name
			// looking like a list so the single-item Get (same ID) isn't picked up.
			if strings.HasPrefix(name, "List") && (res.ListByParentOp == "" || name == "List") {
				res.ListByParentOp = name
				res.ParentIDType = *op.ResourceIDName
			}
		}
	}
	if res.ParentIDType != "" {
		base := strings.TrimSuffix(res.ParentIDType, "Id")
		res.ParentListAttr = snake(base) + "_id"
		res.ParentParseFunc = "Parse" + idToID(res.ParentIDType)
		res.ParentValidateFunc = "Validate" + idToID(res.ParentIDType)
	}
	res.IsListable = res.ListBySubscriptionOp != "" || res.ListByResourceGroupOp != "" || res.ListByParentOp != ""
}

// refName safely extracts a ReferenceName from an ObjectDefinition.
func refName(def *pandora.ObjectDefinition) string {
	if def == nil || def.ReferenceName == nil {
		return ""
	}
	return *def.ReferenceName
}

// deriveID resolves the Go ID type/function names and the ordered arguments for
// the New*ID() constructor from the resource ID segments.
func deriveID(res *ResourceIR, schema *pandora.SchemaResponse, primaryID string) error {
	res.IDTypeName = primaryID
	res.IDParseFunc = "Parse" + idToID(primaryID)
	res.IDNewFunc = "New" + idToID(primaryID)
	res.IDValidateFunc = "Validate" + idToID(primaryID)

	rid, ok := schema.ResourceIds[primaryID]
	if !ok {
		return fmt.Errorf("resource id %q not present in schema", primaryID)
	}

	type userSeg struct {
		argIndex int
		name     string
	}
	var userSegs []userSeg
	args := make([]string, 0, len(rid.Segments))
	for _, seg := range rid.Segments {
		switch seg.Type {
		case "Static", "ResourceProvider":
			continue
		case "SubscriptionId":
			res.HasSubscription = true
			args = append(args, "subscriptionId")
		case "ResourceGroup":
			res.HasResourceGroup = true
			args = append(args, "config.ResourceGroup")
		case "UserSpecified":
			args = append(args, "")
			userSegs = append(userSegs, userSeg{argIndex: len(args) - 1, name: seg.Name})
		default:
			args = append(args, "config."+camel(seg.Name))
		}
	}

	// The final user-specified segment is the resource's own Name; earlier
	// user-specified segments are parent names.
	for i, us := range userSegs {
		if i == len(userSegs)-1 {
			args[us.argIndex] = "config.Name"
			res.IDNameSegment = camel(us.name)
		} else {
			field := camel(strings.TrimSuffix(us.name, "Name"))
			args[us.argIndex] = "config." + field
		}
	}

	res.IDNewArgs = args
	return nil
}
