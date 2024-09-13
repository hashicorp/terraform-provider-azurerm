// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terraform

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-uuid"
	"github.com/mitchellh/copystructure"

	"github.com/hashicorp/terraform-plugin-sdk/v2/internal/addrs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/internal/configs/hcl2shim"
)

const (
	// StateVersion is the current version for our state file
	stateVersion = 3
)

// rootModulePath is the path of the root module
var rootModulePath = []string{"root"}

// normalizeModulePath transforms a legacy module path (which may or may not
// have a redundant "root" label at the start of it) into an
// addrs.ModuleInstance representing the same module.
//
// For legacy reasons, different parts of Terraform disagree about whether the
// root module has the path []string{} or []string{"root"}, and so this
// function accepts both and trims off the "root". An implication of this is
// that it's not possible to actually have a module call in the root module
// that is itself named "root", since that would be ambiguous.
//
// normalizeModulePath takes a raw module path and returns a path that
// has the rootModulePath prepended to it. If I could go back in time I
// would've never had a rootModulePath (empty path would be root). We can
// still fix this but thats a big refactor that my branch doesn't make sense
// for. Instead, this function normalizes paths.
func normalizeModulePath(p []string) addrs.ModuleInstance {
	// FIXME: Remove this once everyone is using addrs.ModuleInstance.

	if len(p) > 0 && p[0] == "root" {
		p = p[1:]
	}

	ret := make(addrs.ModuleInstance, len(p))
	for i, name := range p {
		// For now we don't actually support modules with multiple instances
		// identified by keys, so we just treat every path element as a
		// step with no key.
		ret[i] = addrs.ModuleInstanceStep{
			Name: name,
		}
	}
	return ret
}

// State keeps track of a snapshot state-of-the-world that Terraform
// can use to keep track of what real world resources it is actually
// managing.
type State struct {
	// Version is the state file protocol version.
	Version int `json:"version"`

	// TFVersion is the version of Terraform that wrote this state.
	TFVersion string `json:"terraform_version,omitempty"`

	// Serial is incremented on any operation that modifies
	// the State file. It is used to detect potentially conflicting
	// updates.
	Serial int64 `json:"serial"`

	// Lineage is set when a new, blank state is created and then
	// never updated. This allows us to determine whether the serials
	// of two states can be meaningfully compared.
	// Apart from the guarantee that collisions between two lineages
	// are very unlikely, this value is opaque and external callers
	// should only compare lineage strings byte-for-byte for equality.
	Lineage string `json:"lineage"`

	// Remote is used to track the metadata required to
	// pull and push state files from a remote storage endpoint.
	Remote *RemoteState `json:"remote,omitempty"`

	// Backend tracks the configuration for the backend in use with
	// this state. This is used to track any changes in the backend
	// configuration.
	Backend *BackendState `json:"backend,omitempty"`

	// Modules contains all the modules in a breadth-first order
	Modules []*ModuleState `json:"modules"`

	mu sync.Mutex

	// IsBinaryDrivenTest is a special flag that assists with a binary driver
	// heuristic, it should not be set externally
	IsBinaryDrivenTest bool
}

func (s *State) Lock()   { s.mu.Lock() }
func (s *State) Unlock() { s.mu.Unlock() }

// NewState is used to initialize a blank state
func NewState() *State {
	s := &State{}
	s.init()
	return s
}

// Children returns the ModuleStates that are direct children of
// the given path. If the path is "root", for example, then children
// returned might be "root.child", but not "root.child.grandchild".
func (s *State) Children(path []string) []*ModuleState {
	s.Lock()
	defer s.Unlock()
	// TODO: test

	return s.children(path)
}

func (s *State) children(path []string) []*ModuleState {
	result := make([]*ModuleState, 0)
	for _, m := range s.Modules {
		if m == nil {
			continue
		}

		if len(m.Path) != len(path)+1 {
			continue
		}
		if !reflect.DeepEqual(path, m.Path[:len(path)]) {
			continue
		}

		result = append(result, m)
	}

	return result
}

// AddModule adds the module with the given path to the state.
//
// This should be the preferred method to add module states since it
// allows us to optimize lookups later as well as control sorting.
func (s *State) AddModule(path addrs.ModuleInstance) *ModuleState {
	s.Lock()
	defer s.Unlock()

	return s.addModule(path)
}

func (s *State) addModule(path addrs.ModuleInstance) *ModuleState {
	// check if the module exists first
	m := s.moduleByPath(path)
	if m != nil {
		return m
	}

	// Lower the new-style address into a legacy-style address.
	// This requires that none of the steps have instance keys, which is
	// true for all addresses at the time of implementing this because
	// "count" and "for_each" are not yet implemented for modules.
	// For the purposes of state, the legacy address format also includes
	// a redundant extra prefix element "root". It is important to include
	// this because the "prune" method will remove any module that has a
	// path length less than one, and other parts of the state code will
	// trim off the first element indiscriminately.
	legacyPath := make([]string, len(path)+1)
	legacyPath[0] = "root"
	for i, step := range path {
		if step.InstanceKey != addrs.NoKey {
			// FIXME: Once the rest of Terraform is ready to use count and
			// for_each, remove all of this and just write the addrs.ModuleInstance
			// value itself into the ModuleState.
			panic("state cannot represent modules with count or for_each keys")
		}

		legacyPath[i+1] = step.Name
	}

	m = &ModuleState{Path: legacyPath}
	m.init()
	s.Modules = append(s.Modules, m)
	s.sort()
	return m
}

// ModuleByPath is used to lookup the module state for the given path.
// This should be the preferred lookup mechanism as it allows for future
// lookup optimizations.
func (s *State) ModuleByPath(path addrs.ModuleInstance) *ModuleState {
	if s == nil {
		return nil
	}
	s.Lock()
	defer s.Unlock()

	return s.moduleByPath(path)
}

func (s *State) moduleByPath(path addrs.ModuleInstance) *ModuleState {
	for _, mod := range s.Modules {
		if mod == nil {
			continue
		}
		if mod.Path == nil {
			panic("missing module path")
		}
		modPath := normalizeModulePath(mod.Path)
		if modPath.String() == path.String() {
			return mod
		}
	}
	return nil
}

// Empty returns true if the state is empty.
func (s *State) Empty() bool {
	if s == nil {
		return true
	}
	s.Lock()
	defer s.Unlock()

	return len(s.Modules) == 0
}

// HasResources returns true if the state contains any resources.
//
// This is similar to !s.Empty, but returns true also in the case where the
// state has modules but all of them are devoid of resources.
func (s *State) HasResources() bool {
	if s.Empty() {
		return false
	}

	for _, mod := range s.Modules {
		if len(mod.Resources) > 0 {
			return true
		}
	}

	return false
}

// IsRemote returns true if State represents a state that exists and is
// remote.
func (s *State) IsRemote() bool {
	if s == nil {
		return false
	}
	s.Lock()
	defer s.Unlock()

	if s.Remote == nil {
		return false
	}
	if s.Remote.Type == "" {
		return false
	}

	return true
}

// Validate validates the integrity of this state file.
//
// Certain properties of the statefile are expected by Terraform in order
// to behave properly. The core of Terraform will assume that once it
// receives a State structure that it has been validated. This validation
// check should be called to ensure that.
//
// If this returns an error, then the user should be notified. The error
// response will include detailed information on the nature of the error.
func (s *State) Validate() error {
	s.Lock()
	defer s.Unlock()

	var result []error

	// !!!! FOR DEVELOPERS !!!!
	//
	// Any errors returned from this Validate function will BLOCK TERRAFORM
	// from loading a state file. Therefore, this should only contain checks
	// that are only resolvable through manual intervention.
	//
	// !!!! FOR DEVELOPERS !!!!

	// Make sure there are no duplicate module states. We open a new
	// block here so we can use basic variable names and future validations
	// can do the same.
	{
		found := make(map[string]struct{})
		for _, ms := range s.Modules {
			if ms == nil {
				continue
			}

			key := strings.Join(ms.Path, ".")
			if _, ok := found[key]; ok {
				result = append(result, fmt.Errorf(
					strings.TrimSpace(stateValidateErrMultiModule), key))
				continue
			}

			found[key] = struct{}{}
		}
	}

	return errors.Join(result...)
}

// Remove removes the item in the state at the given address, returning
// any errors that may have occurred.
//
// If the address references a module state or resource, it will delete
// all children as well. To check what will be deleted, use a StateFilter
// first.
func (s *State) Remove(addr ...string) error {
	s.Lock()
	defer s.Unlock()

	// Filter out what we need to delete
	filter := &stateFilter{State: s}
	results, err := filter.filter(addr...)
	if err != nil {
		return err
	}

	// If we have no results, just exit early, we're not going to do anything.
	// While what happens below is fairly fast, this is an important early
	// exit since the prune below might modify the state more and we don't
	// want to modify the state if we don't have to.
	if len(results) == 0 {
		return nil
	}

	// Go through each result and grab what we need
	removed := make(map[interface{}]struct{})
	for _, r := range results {
		// Convert the path to our own type
		path := append([]string{"root"}, r.Path...)

		// If we removed this already, then ignore
		if _, ok := removed[r.Value]; ok {
			continue
		}

		// If we removed the parent already, then ignore
		if r.Parent != nil {
			if _, ok := removed[r.Parent.Value]; ok {
				continue
			}
		}

		// Add this to the removed list
		removed[r.Value] = struct{}{}

		switch v := r.Value.(type) {
		case *ModuleState:
			s.removeModule(v)
		case *ResourceState:
			s.removeResource(path, v)
		case *InstanceState:
			s.removeInstance(r.Parent.Value.(*ResourceState), v)
		default:
			return fmt.Errorf("unknown type to delete: %T", r.Value)
		}
	}

	// Prune since the removal functions often do the bare minimum to
	// remove a thing and may leave around dangling empty modules, resources,
	// etc. Prune will clean that all up.
	s.prune()

	return nil
}

func (s *State) removeModule(v *ModuleState) {
	for i, m := range s.Modules {
		if m == v {
			s.Modules, s.Modules[len(s.Modules)-1] = append(s.Modules[:i], s.Modules[i+1:]...), nil
			return
		}
	}
}

func (s *State) removeResource(path []string, v *ResourceState) {
	// Get the module this resource lives in. If it doesn't exist, we're done.
	mod := s.moduleByPath(normalizeModulePath(path))
	if mod == nil {
		return
	}

	// Find this resource. This is a O(N) lookup when if we had the key
	// it could be O(1) but even with thousands of resources this shouldn't
	// matter right now. We can easily up performance here when the time comes.
	for k, r := range mod.Resources {
		if r == v {
			// Found it
			delete(mod.Resources, k)
			return
		}
	}
}

func (s *State) removeInstance(r *ResourceState, v *InstanceState) {
	// Go through the resource and find the instance that matches this
	// (if any) and remove it.

	// Check primary
	if r.Primary == v {
		r.Primary = nil
		return
	}
}

// RootModule returns the ModuleState for the root module
func (s *State) RootModule() *ModuleState {
	root := s.ModuleByPath(addrs.RootModuleInstance)
	if root == nil {
		panic("missing root module")
	}
	return root
}

// Equal tests if one state is equal to another.
func (s *State) Equal(other *State) bool {
	// If one is nil, we do a direct check
	if s == nil || other == nil {
		return s == other
	}

	s.Lock()
	defer s.Unlock()
	return s.equal(other)
}

func (s *State) equal(other *State) bool {
	if s == nil || other == nil {
		return s == other
	}

	// If the versions are different, they're certainly not equal
	if s.Version != other.Version {
		return false
	}

	// If any of the modules are not equal, then this state isn't equal
	if len(s.Modules) != len(other.Modules) {
		return false
	}
	for _, m := range s.Modules {
		// This isn't very optimal currently but works.
		otherM := other.moduleByPath(normalizeModulePath(m.Path))
		if otherM == nil {
			return false
		}

		// If they're not equal, then we're not equal!
		if !m.Equal(otherM) {
			return false
		}
	}

	return true
}

type StateAgeComparison int

const (
	StateAgeEqual         StateAgeComparison = 0
	StateAgeReceiverNewer StateAgeComparison = 1
	StateAgeReceiverOlder StateAgeComparison = -1
)

// CompareAges compares one state with another for which is "older".
//
// This is a simple check using the state's serial, and is thus only as
// reliable as the serial itself. In the normal case, only one state
// exists for a given combination of lineage/serial, but Terraform
// does not guarantee this and so the result of this method should be
// used with care.
//
// Returns an integer that is negative if the receiver is older than
// the argument, positive if the converse, and zero if they are equal.
// An error is returned if the two states are not of the same lineage,
// in which case the integer returned has no meaning.
func (s *State) CompareAges(other *State) (StateAgeComparison, error) {
	// nil states are "older" than actual states
	switch {
	case s != nil && other == nil:
		return StateAgeReceiverNewer, nil
	case s == nil && other != nil:
		return StateAgeReceiverOlder, nil
	case s == nil && other == nil:
		return StateAgeEqual, nil
	}

	if !s.SameLineage(other) {
		return StateAgeEqual, fmt.Errorf(
			"can't compare two states of differing lineage",
		)
	}

	s.Lock()
	defer s.Unlock()

	switch {
	case s.Serial < other.Serial:
		return StateAgeReceiverOlder, nil
	case s.Serial > other.Serial:
		return StateAgeReceiverNewer, nil
	default:
		return StateAgeEqual, nil
	}
}

// SameLineage returns true only if the state given in argument belongs
// to the same "lineage" of states as the receiver.
func (s *State) SameLineage(other *State) bool {
	s.Lock()
	defer s.Unlock()

	// If one of the states has no lineage then it is assumed to predate
	// this concept, and so we'll accept it as belonging to any lineage
	// so that a lineage string can be assigned to newer versions
	// without breaking compatibility with older versions.
	if s.Lineage == "" || other.Lineage == "" {
		return true
	}

	return s.Lineage == other.Lineage
}

// DeepCopy performs a deep copy of the state structure and returns
// a new structure.
func (s *State) DeepCopy() *State {
	if s == nil {
		return nil
	}

	copiedState, err := copystructure.Config{Lock: true}.Copy(s)
	if err != nil {
		panic(err)
	}

	return copiedState.(*State)
}

func (s *State) Init() {
	s.Lock()
	defer s.Unlock()
	s.init()
}

func (s *State) init() {
	if s.Version == 0 {
		s.Version = stateVersion
	}

	if s.moduleByPath(addrs.RootModuleInstance) == nil {
		s.addModule(addrs.RootModuleInstance)
	}
	s.ensureHasLineage()

	for _, mod := range s.Modules {
		if mod != nil {
			mod.init()
		}
	}

	if s.Remote != nil {
		s.Remote.init()
	}

}

func (s *State) EnsureHasLineage() {
	s.Lock()
	defer s.Unlock()

	s.ensureHasLineage()
}

func (s *State) ensureHasLineage() {
	if s.Lineage == "" {
		lineage, err := uuid.GenerateUUID()
		if err != nil {
			panic(fmt.Errorf("Failed to generate lineage: %v", err))
		}
		s.Lineage = lineage
		if os.Getenv("TF_ACC") == "" || os.Getenv("TF_ACC_STATE_LINEAGE") == "1" {
			log.Printf("[DEBUG] New state was assigned lineage %q\n", s.Lineage)
		}
	} else {
		if os.Getenv("TF_ACC") == "" || os.Getenv("TF_ACC_STATE_LINEAGE") == "1" {
			log.Printf("[TRACE] Preserving existing state lineage %q\n", s.Lineage)
		}
	}
}

// AddModuleState insert this module state and override any existing ModuleState
func (s *State) AddModuleState(mod *ModuleState) {
	mod.init()
	s.Lock()
	defer s.Unlock()

	s.addModuleState(mod)
}

func (s *State) addModuleState(mod *ModuleState) {
	for i, m := range s.Modules {
		if reflect.DeepEqual(m.Path, mod.Path) {
			s.Modules[i] = mod
			return
		}
	}

	s.Modules = append(s.Modules, mod)
	s.sort()
}

// prune is used to remove any resources that are no longer required
func (s *State) prune() {
	if s == nil {
		return
	}

	// Filter out empty modules.
	// A module is always assumed to have a path, and it's length isn't always
	// bounds checked later on. Modules may be "emptied" during destroy, but we
	// never want to store those in the state.
	for i := 0; i < len(s.Modules); i++ {
		if s.Modules[i] == nil || len(s.Modules[i].Path) == 0 {
			s.Modules = append(s.Modules[:i], s.Modules[i+1:]...)
			i--
		}
	}

	for _, mod := range s.Modules {
		mod.prune()
	}
	if s.Remote != nil && s.Remote.Empty() {
		s.Remote = nil
	}
}

// sort sorts the modules
func (s *State) sort() {
	sort.Sort(moduleStateSort(s.Modules))

	// Allow modules to be sorted
	for _, m := range s.Modules {
		if m != nil {
			m.sort()
		}
	}
}

func (s *State) String() string {
	if s == nil {
		return "<nil>"
	}
	s.Lock()
	defer s.Unlock()

	var buf bytes.Buffer
	for _, m := range s.Modules {
		mStr := m.String()

		// If we're the root module, we just write the output directly.
		if reflect.DeepEqual(m.Path, rootModulePath) {
			buf.WriteString(mStr + "\n")
			continue
		}

		buf.WriteString(fmt.Sprintf("module.%s:\n", strings.Join(m.Path[1:], ".")))

		s := bufio.NewScanner(strings.NewReader(mStr))
		for s.Scan() {
			text := s.Text()
			if text != "" {
				text = "  " + text
			}

			buf.WriteString(fmt.Sprintf("%s\n", text))
		}
	}

	return strings.TrimSpace(buf.String())
}

// BackendState stores the configuration to connect to a remote backend.
type BackendState struct {
	Type      string          `json:"type"`   // Backend type
	ConfigRaw json.RawMessage `json:"config"` // Backend raw config
	Hash      uint64          `json:"hash"`   // Hash of portion of configuration from config files
}

// RemoteState is used to track the information about a remote
// state store that we push/pull state to.
type RemoteState struct {
	// Type controls the client we use for the remote state
	Type string `json:"type"`

	// Config is used to store arbitrary configuration that
	// is type specific
	Config map[string]string `json:"config"`

	mu sync.Mutex
}

func (s *RemoteState) Lock()   { s.mu.Lock() }
func (s *RemoteState) Unlock() { s.mu.Unlock() }

func (r *RemoteState) init() {
	r.Lock()
	defer r.Unlock()

	if r.Config == nil {
		r.Config = make(map[string]string)
	}
}

func (r *RemoteState) Empty() bool {
	if r == nil {
		return true
	}
	r.Lock()
	defer r.Unlock()

	return r.Type == ""
}

// OutputState is used to track the state relevant to a single output.
type OutputState struct {
	// Sensitive describes whether the output is considered sensitive,
	// which may lead to masking the value on screen in some cases.
	Sensitive bool `json:"sensitive"`
	// Type describes the structure of Value. Valid values are "string",
	// "map" and "list"
	Type string `json:"type"`
	// Value contains the value of the output, in the structure described
	// by the Type field.
	Value interface{} `json:"value"`

	mu sync.Mutex
}

func (s *OutputState) Lock()   { s.mu.Lock() }
func (s *OutputState) Unlock() { s.mu.Unlock() }

func (s *OutputState) String() string {
	return fmt.Sprintf("%#v", s.Value)
}

// Equal compares two OutputState structures for equality. nil values are
// considered equal.
func (s *OutputState) Equal(other *OutputState) bool {
	if s == nil && other == nil {
		return true
	}

	if s == nil || other == nil {
		return false
	}
	s.Lock()
	defer s.Unlock()

	if s.Type != other.Type {
		return false
	}

	if s.Sensitive != other.Sensitive {
		return false
	}

	if !reflect.DeepEqual(s.Value, other.Value) {
		return false
	}

	return true
}

// ModuleState is used to track all the state relevant to a single
// module. Previous to Terraform 0.3, all state belonged to the "root"
// module.
type ModuleState struct {
	// Path is the import path from the root module. Modules imports are
	// always disjoint, so the path represents amodule tree
	Path []string `json:"path"`

	// Locals are kept only transiently in-memory, because we can always
	// re-compute them.
	Locals map[string]interface{} `json:"-"`

	// Outputs declared by the module and maintained for each module
	// even though only the root module technically needs to be kept.
	// This allows operators to inspect values at the boundaries.
	Outputs map[string]*OutputState `json:"outputs"`

	// Resources is a mapping of the logically named resource to
	// the state of the resource. Each resource may actually have
	// N instances underneath, although a user only needs to think
	// about the 1:1 case.
	Resources map[string]*ResourceState `json:"resources"`

	// Dependencies are a list of things that this module relies on
	// existing to remain intact. For example: an module may depend
	// on a VPC ID given by an aws_vpc resource.
	//
	// Terraform uses this information to build valid destruction
	// orders and to warn the user if they're destroying a module that
	// another resource depends on.
	//
	// Things can be put into this list that may not be managed by
	// Terraform. If Terraform doesn't find a matching ID in the
	// overall state, then it assumes it isn't managed and doesn't
	// worry about it.
	Dependencies []string `json:"depends_on"`

	mu sync.Mutex
}

func (s *ModuleState) Lock()   { s.mu.Lock() }
func (s *ModuleState) Unlock() { s.mu.Unlock() }

// Equal tests whether one module state is equal to another.
func (m *ModuleState) Equal(other *ModuleState) bool {
	m.Lock()
	defer m.Unlock()

	// Paths must be equal
	if !reflect.DeepEqual(m.Path, other.Path) {
		return false
	}

	// Outputs must be equal
	if len(m.Outputs) != len(other.Outputs) {
		return false
	}
	for k, v := range m.Outputs {
		if !other.Outputs[k].Equal(v) {
			return false
		}
	}

	// Dependencies must be equal. This sorts these in place but
	// this shouldn't cause any problems.
	sort.Strings(m.Dependencies)
	sort.Strings(other.Dependencies)
	if len(m.Dependencies) != len(other.Dependencies) {
		return false
	}
	for i, d := range m.Dependencies {
		if other.Dependencies[i] != d {
			return false
		}
	}

	// Resources must be equal
	if len(m.Resources) != len(other.Resources) {
		return false
	}
	for k, r := range m.Resources {
		otherR, ok := other.Resources[k]
		if !ok {
			return false
		}

		if !r.Equal(otherR) {
			return false
		}
	}

	return true
}

func (m *ModuleState) init() {
	m.Lock()
	defer m.Unlock()

	if m.Path == nil {
		m.Path = []string{}
	}
	if m.Outputs == nil {
		m.Outputs = make(map[string]*OutputState)
	}
	if m.Resources == nil {
		m.Resources = make(map[string]*ResourceState)
	}

	if m.Dependencies == nil {
		m.Dependencies = make([]string, 0)
	}

	for _, rs := range m.Resources {
		rs.init()
	}
}

// prune is used to remove any resources that are no longer required
func (m *ModuleState) prune() {
	m.Lock()
	defer m.Unlock()

	for k, v := range m.Resources {
		if v == nil || (v.Primary == nil || v.Primary.ID == "") && len(v.Deposed) == 0 {
			delete(m.Resources, k)
			continue
		}

		v.prune()
	}

	for k, v := range m.Outputs {
		if v.Value == hcl2shim.UnknownVariableValue {
			delete(m.Outputs, k)
		}
	}

	m.Dependencies = uniqueStrings(m.Dependencies)
}

func (m *ModuleState) sort() {
	for _, v := range m.Resources {
		v.sort()
	}
}

func (m *ModuleState) String() string {
	m.Lock()
	defer m.Unlock()

	var buf bytes.Buffer

	if len(m.Resources) == 0 {
		buf.WriteString("<no state>")
	}

	names := make([]string, 0, len(m.Resources))
	for name := range m.Resources {
		names = append(names, name)
	}

	sort.Sort(resourceNameSort(names))

	for _, k := range names {
		rs := m.Resources[k]
		var id string
		if rs.Primary != nil {
			id = rs.Primary.ID
		}
		if id == "" {
			id = "<not created>"
		}

		taintStr := ""
		if rs.Primary.Tainted {
			taintStr = " (tainted)"
		}

		deposedStr := ""
		if len(rs.Deposed) > 0 {
			deposedStr = fmt.Sprintf(" (%d deposed)", len(rs.Deposed))
		}

		buf.WriteString(fmt.Sprintf("%s:%s%s\n", k, taintStr, deposedStr))
		buf.WriteString(fmt.Sprintf("  ID = %s\n", id))
		if rs.Provider != "" {
			buf.WriteString(fmt.Sprintf("  provider = %s\n", rs.Provider))
		}

		var attributes map[string]string
		if rs.Primary != nil {
			attributes = rs.Primary.Attributes
		}
		attrKeys := make([]string, 0, len(attributes))
		for ak := range attributes {
			if ak == "id" {
				continue
			}

			attrKeys = append(attrKeys, ak)
		}

		sort.Strings(attrKeys)

		for _, ak := range attrKeys {
			av := attributes[ak]
			buf.WriteString(fmt.Sprintf("  %s = %s\n", ak, av))
		}

		for idx, t := range rs.Deposed {
			taintStr := ""
			if t.Tainted {
				taintStr = " (tainted)"
			}
			buf.WriteString(fmt.Sprintf("  Deposed ID %d = %s%s\n", idx+1, t.ID, taintStr))
		}

		if len(rs.Dependencies) > 0 {
			buf.WriteString("\n  Dependencies:\n")
			for _, dep := range rs.Dependencies {
				buf.WriteString(fmt.Sprintf("    %s\n", dep))
			}
		}
	}

	if len(m.Outputs) > 0 {
		buf.WriteString("\nOutputs:\n\n")

		ks := make([]string, 0, len(m.Outputs))
		for k := range m.Outputs {
			ks = append(ks, k)
		}

		sort.Strings(ks)

		for _, k := range ks {
			v := m.Outputs[k]
			switch vTyped := v.Value.(type) {
			case string:
				buf.WriteString(fmt.Sprintf("%s = %s\n", k, vTyped))
			case []interface{}:
				buf.WriteString(fmt.Sprintf("%s = %s\n", k, vTyped))
			case map[string]interface{}:
				var mapKeys []string
				for key := range vTyped {
					mapKeys = append(mapKeys, key)
				}
				sort.Strings(mapKeys)

				var mapBuf bytes.Buffer
				mapBuf.WriteString("{")
				for _, key := range mapKeys {
					mapBuf.WriteString(fmt.Sprintf("%s:%s ", key, vTyped[key]))
				}
				mapBuf.WriteString("}")

				buf.WriteString(fmt.Sprintf("%s = %s\n", k, mapBuf.String()))
			}
		}
	}

	return buf.String()
}

// ResourceStateKey is a structured representation of the key used for the
// ModuleState.Resources mapping
type ResourceStateKey struct {
	Name  string
	Type  string
	Mode  ResourceMode
	Index int
}

// Equal determines whether two ResourceStateKeys are the same
func (rsk *ResourceStateKey) Equal(other *ResourceStateKey) bool {
	if rsk == nil || other == nil {
		return false
	}
	if rsk.Mode != other.Mode {
		return false
	}
	if rsk.Type != other.Type {
		return false
	}
	if rsk.Name != other.Name {
		return false
	}
	if rsk.Index != other.Index {
		return false
	}
	return true
}

func (rsk *ResourceStateKey) String() string {
	if rsk == nil {
		return ""
	}
	var prefix string
	switch rsk.Mode {
	case ManagedResourceMode:
		prefix = ""
	case DataResourceMode:
		prefix = "data."
	default:
		panic(fmt.Errorf("unknown resource mode %s", rsk.Mode))
	}
	if rsk.Index == -1 {
		return fmt.Sprintf("%s%s.%s", prefix, rsk.Type, rsk.Name)
	}
	return fmt.Sprintf("%s%s.%s.%d", prefix, rsk.Type, rsk.Name, rsk.Index)
}

// ParseResourceStateKey accepts a key in the format used by
// ModuleState.Resources and returns a resource name and resource index. In the
// state, a resource has the format "type.name.index" or "type.name". In the
// latter case, the index is returned as -1.
func parseResourceStateKey(k string) (*ResourceStateKey, error) {
	parts := strings.Split(k, ".")
	mode := ManagedResourceMode
	if len(parts) > 0 && parts[0] == "data" {
		mode = DataResourceMode
		// Don't need the constant "data" prefix for parsing
		// now that we've figured out the mode.
		parts = parts[1:]
	}
	if len(parts) < 2 || len(parts) > 3 {
		return nil, fmt.Errorf("Malformed resource state key: %s", k)
	}
	rsk := &ResourceStateKey{
		Mode:  mode,
		Type:  parts[0],
		Name:  parts[1],
		Index: -1,
	}
	if len(parts) == 3 {
		index, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("Malformed resource state key index: %s", k)
		}
		rsk.Index = index
	}
	return rsk, nil
}

// ResourceState holds the state of a resource that is used so that
// a provider can find and manage an existing resource as well as for
// storing attributes that are used to populate variables of child
// resources.
//
// Attributes has attributes about the created resource that are
// queryable in interpolation: "${type.id.attr}"
//
// Extra is just extra data that a provider can return that we store
// for later, but is not exposed in any way to the user.
type ResourceState struct {
	// This is filled in and managed by Terraform, and is the resource
	// type itself such as "mycloud_instance". If a resource provider sets
	// this value, it won't be persisted.
	Type string `json:"type"`

	// Dependencies are a list of things that this resource relies on
	// existing to remain intact. For example: an AWS instance might
	// depend on a subnet (which itself might depend on a VPC, and so
	// on).
	//
	// Terraform uses this information to build valid destruction
	// orders and to warn the user if they're destroying a resource that
	// another resource depends on.
	//
	// Things can be put into this list that may not be managed by
	// Terraform. If Terraform doesn't find a matching ID in the
	// overall state, then it assumes it isn't managed and doesn't
	// worry about it.
	Dependencies []string `json:"depends_on"`

	// Primary is the current active instance for this resource.
	// It can be replaced but only after a successful creation.
	// This is the instances on which providers will act.
	Primary *InstanceState `json:"primary"`

	// Deposed is used in the mechanics of CreateBeforeDestroy: the existing
	// Primary is Deposed to get it out of the way for the replacement Primary to
	// be created by Apply. If the replacement Primary creates successfully, the
	// Deposed instance is cleaned up.
	//
	// If there were problems creating the replacement Primary, the Deposed
	// instance and the (now tainted) replacement Primary will be swapped so the
	// tainted replacement will be cleaned up instead.
	//
	// An instance will remain in the Deposed list until it is successfully
	// destroyed and purged.
	Deposed []*InstanceState `json:"deposed"`

	// Provider is used when a resource is connected to a provider with an alias.
	// If this string is empty, the resource is connected to the default provider,
	// e.g. "aws_instance" goes with the "aws" provider.
	// If the resource block contained a "provider" key, that value will be set here.
	Provider string `json:"provider"`

	mu sync.Mutex
}

func (s *ResourceState) Lock()   { s.mu.Lock() }
func (s *ResourceState) Unlock() { s.mu.Unlock() }

// Equal tests whether two ResourceStates are equal.
func (s *ResourceState) Equal(other *ResourceState) bool {
	s.Lock()
	defer s.Unlock()

	if s.Type != other.Type {
		return false
	}

	if s.Provider != other.Provider {
		return false
	}

	// Dependencies must be equal
	sort.Strings(s.Dependencies)
	sort.Strings(other.Dependencies)
	if len(s.Dependencies) != len(other.Dependencies) {
		return false
	}
	for i, d := range s.Dependencies {
		if other.Dependencies[i] != d {
			return false
		}
	}

	// States must be equal
	return s.Primary.Equal(other.Primary)
}

func (s *ResourceState) init() {
	s.Lock()
	defer s.Unlock()

	if s.Primary == nil {
		s.Primary = &InstanceState{}
	}
	s.Primary.init()

	if s.Dependencies == nil {
		s.Dependencies = []string{}
	}

	if s.Deposed == nil {
		s.Deposed = make([]*InstanceState, 0)
	}
}

// prune is used to remove any instances that are no longer required
func (s *ResourceState) prune() {
	s.Lock()
	defer s.Unlock()

	n := len(s.Deposed)
	for i := 0; i < n; i++ {
		inst := s.Deposed[i]
		if inst == nil || inst.ID == "" {
			copy(s.Deposed[i:], s.Deposed[i+1:])
			s.Deposed[n-1] = nil
			n--
			i--
		}
	}
	s.Deposed = s.Deposed[:n]

	s.Dependencies = uniqueStrings(s.Dependencies)
}

func (s *ResourceState) sort() {
	s.Lock()
	defer s.Unlock()

	sort.Strings(s.Dependencies)
}

func (s *ResourceState) String() string {
	s.Lock()
	defer s.Unlock()

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Type = %s", s.Type))
	return buf.String()
}

// InstanceState is used to track the unique state information belonging
// to a given instance.
type InstanceState struct {
	// A unique ID for this resource. This is opaque to Terraform
	// and is only meant as a lookup mechanism for the providers.
	ID string `json:"id"`

	// Attributes are basic information about the resource. Any keys here
	// are accessible in variable format within Terraform configurations:
	// ${resourcetype.name.attribute}.
	Attributes map[string]string `json:"attributes"`

	// Ephemeral is used to store any state associated with this instance
	// that is necessary for the Terraform run to complete, but is not
	// persisted to a state file.
	Ephemeral EphemeralState `json:"-"`

	// Meta is a simple K/V map that is persisted to the State but otherwise
	// ignored by Terraform core. It's meant to be used for accounting by
	// external client code. The value here must only contain Go primitives
	// and collections.
	Meta map[string]interface{} `json:"meta"`

	ProviderMeta cty.Value

	RawConfig cty.Value
	RawState  cty.Value
	RawPlan   cty.Value

	// Tainted is used to mark a resource for recreation.
	Tainted bool `json:"tainted"`

	mu sync.Mutex
}

func (s *InstanceState) Lock()   { s.mu.Lock() }
func (s *InstanceState) Unlock() { s.mu.Unlock() }

func (s *InstanceState) init() {
	s.Lock()
	defer s.Unlock()

	if s.Attributes == nil {
		s.Attributes = make(map[string]string)
	}
	if s.Meta == nil {
		s.Meta = make(map[string]interface{})
	}
	s.Ephemeral.init()
}

// NewInstanceStateShimmedFromValue is a shim method to lower a new-style
// object value representing the attributes of an instance object into the
// legacy InstanceState representation.
//
// This is for shimming to old components only and should not be used in new code.
func NewInstanceStateShimmedFromValue(state cty.Value, schemaVersion int) *InstanceState {
	attrs := hcl2shim.FlatmapValueFromHCL2(state)
	return &InstanceState{
		ID:         attrs["id"],
		Attributes: attrs,
		Meta: map[string]interface{}{
			"schema_version": schemaVersion,
		},
	}
}

// AttrsAsObjectValue shims from the legacy InstanceState representation to
// a new-style cty object value representation of the state attributes, using
// the given type for guidance.
//
// The given type must be the implied type of the schema of the resource type
// of the object whose state is being converted, or the result is undefined.
//
// This is for shimming from old components only and should not be used in
// new code.
func (s *InstanceState) AttrsAsObjectValue(ty cty.Type) (cty.Value, error) {
	if s == nil {
		// if the state is nil, we need to construct a complete cty.Value with
		// null attributes, rather than a single cty.NullVal(ty)
		s = &InstanceState{}
	}

	if s.Attributes == nil {
		s.Attributes = map[string]string{}
	}

	// make sure ID is included in the attributes. The InstanceState.ID value
	// takes precedence.
	if s.ID != "" {
		s.Attributes["id"] = s.ID
	}

	return hcl2shim.HCL2ValueFromFlatmap(s.Attributes, ty)
}

// Copy all the Fields from another InstanceState
func (s *InstanceState) Set(from *InstanceState) {
	s.Lock()
	defer s.Unlock()

	from.Lock()
	defer from.Unlock()

	s.ID = from.ID
	s.Attributes = from.Attributes
	s.Ephemeral = from.Ephemeral
	s.Meta = from.Meta
	s.Tainted = from.Tainted
}

func (s *InstanceState) DeepCopy() *InstanceState {
	copiedState, err := copystructure.Config{Lock: true}.Copy(s)
	if err != nil {
		panic(err)
	}

	return copiedState.(*InstanceState)
}

func (s *InstanceState) Empty() bool {
	if s == nil {
		return true
	}
	s.Lock()
	defer s.Unlock()

	return s.ID == ""
}

func (s *InstanceState) Equal(other *InstanceState) bool {
	// Short circuit some nil checks
	if s == nil || other == nil {
		return s == other
	}
	s.Lock()
	defer s.Unlock()

	// IDs must be equal
	if s.ID != other.ID {
		return false
	}

	// Attributes must be equal
	if len(s.Attributes) != len(other.Attributes) {
		return false
	}
	for k, v := range s.Attributes {
		otherV, ok := other.Attributes[k]
		if !ok {
			return false
		}

		if v != otherV {
			return false
		}
	}

	// Meta must be equal
	if len(s.Meta) != len(other.Meta) {
		return false
	}
	if s.Meta != nil && other.Meta != nil {
		// We only do the deep check if both are non-nil. If one is nil
		// we treat it as equal since their lengths are both zero (check
		// above).
		//
		// Since this can contain numeric values that may change types during
		// serialization, let's compare the serialized values.
		sMeta, err := json.Marshal(s.Meta)
		if err != nil {
			// marshaling primitives shouldn't ever error out
			panic(err)
		}
		otherMeta, err := json.Marshal(other.Meta)
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(sMeta, otherMeta) {
			return false
		}
	}

	if s.Tainted != other.Tainted {
		return false
	}

	return true
}

// MergeDiff takes a ResourceDiff and merges the attributes into
// this resource state in order to generate a new state. This new
// state can be used to provide updated attribute lookups for
// variable interpolation.
//
// If the diff attribute requires computing the value, and hence
// won't be available until apply, the value is replaced with the
// computeID.
func (s *InstanceState) MergeDiff(d *InstanceDiff) *InstanceState {
	result := s.DeepCopy()
	if result == nil {
		result = new(InstanceState)
	}
	result.init()

	if s != nil {
		s.Lock()
		defer s.Unlock()
		for k, v := range s.Attributes {
			result.Attributes[k] = v
		}
	}
	if d != nil {
		for k, diff := range d.CopyAttributes() {
			if diff.NewRemoved {
				delete(result.Attributes, k)
				continue
			}
			if diff.NewComputed {
				result.Attributes[k] = hcl2shim.UnknownVariableValue
				continue
			}

			result.Attributes[k] = diff.New
		}
	}

	return result
}

func (s *InstanceState) String() string {
	notCreated := "<not created>"

	if s == nil {
		return notCreated
	}

	s.Lock()
	defer s.Unlock()

	var buf bytes.Buffer

	if s.ID == "" {
		return notCreated
	}

	buf.WriteString(fmt.Sprintf("ID = %s\n", s.ID))

	attributes := s.Attributes
	attrKeys := make([]string, 0, len(attributes))
	for ak := range attributes {
		if ak == "id" {
			continue
		}

		attrKeys = append(attrKeys, ak)
	}
	sort.Strings(attrKeys)

	for _, ak := range attrKeys {
		av := attributes[ak]
		buf.WriteString(fmt.Sprintf("%s = %s\n", ak, av))
	}

	buf.WriteString(fmt.Sprintf("Tainted = %t\n", s.Tainted))

	return buf.String()
}

// EphemeralState is used for transient state that is only kept in-memory
type EphemeralState struct {
	// ConnInfo is used for the providers to export information which is
	// used to connect to the resource for provisioning. For example,
	// this could contain SSH or WinRM credentials.
	ConnInfo map[string]string `json:"-"`

	// Type is used to specify the resource type for this instance. This is only
	// required for import operations (as documented). If the documentation
	// doesn't state that you need to set this, then don't worry about
	// setting it.
	Type string `json:"-"`
}

func (e *EphemeralState) init() {
	if e.ConnInfo == nil {
		e.ConnInfo = make(map[string]string)
	}
}

// resourceNameSort implements the sort.Interface to sort name parts lexically for
// strings and numerically for integer indexes.
type resourceNameSort []string

func (r resourceNameSort) Len() int      { return len(r) }
func (r resourceNameSort) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func (r resourceNameSort) Less(i, j int) bool {
	iParts := strings.Split(r[i], ".")
	jParts := strings.Split(r[j], ".")

	end := len(iParts)
	if len(jParts) < end {
		end = len(jParts)
	}

	for idx := 0; idx < end; idx++ {
		if iParts[idx] == jParts[idx] {
			continue
		}

		// sort on the first non-matching part
		iInt, iIntErr := strconv.Atoi(iParts[idx])
		jInt, jIntErr := strconv.Atoi(jParts[idx])

		switch {
		case iIntErr == nil && jIntErr == nil:
			// sort numerically if both parts are integers
			return iInt < jInt
		case iIntErr == nil:
			// numbers sort before strings
			return true
		case jIntErr == nil:
			return false
		default:
			return iParts[idx] < jParts[idx]
		}
	}

	return r[i] < r[j]
}

// moduleStateSort implements sort.Interface to sort module states
type moduleStateSort []*ModuleState

func (s moduleStateSort) Len() int {
	return len(s)
}

func (s moduleStateSort) Less(i, j int) bool {
	a := s[i]
	b := s[j]

	// If either is nil, then the nil one is "less" than
	if a == nil || b == nil {
		return a == nil
	}

	// If the lengths are different, then the shorter one always wins
	if len(a.Path) != len(b.Path) {
		return len(a.Path) < len(b.Path)
	}

	// Otherwise, compare lexically
	return strings.Join(a.Path, ".") < strings.Join(b.Path, ".")
}

func (s moduleStateSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

const stateValidateErrMultiModule = `
Multiple modules with the same path: %s

This means that there are multiple entries in the "modules" field
in your state file that point to the same module. This will cause Terraform
to behave in unexpected and error prone ways and is invalid. Please back up
and modify your state file manually to resolve this.
`
