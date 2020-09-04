package tfexec

type AllowMissingConfigOption struct {
	allowMissingConfig bool
}

func AllowMissingConfig(allowMissingConfig bool) *AllowMissingConfigOption {
	return &AllowMissingConfigOption{allowMissingConfig}
}

type BackendOption struct {
	backend bool
}

func Backend(backend bool) *BackendOption {
	return &BackendOption{backend}
}

type BackendConfigOption struct {
	path string
}

func BackendConfig(backendConfig string) *BackendConfigOption {
	return &BackendConfigOption{backendConfig}
}

type BackupOption struct {
	path string
}

func Backup(path string) *BackupOption {
	return &BackupOption{path}
}

// DisableBackup is a convenience method for Backup("-"), indicating backup state should be disabled.
func DisableBackup() *BackupOption {
	return &BackupOption{"-"}
}

type ConfigOption struct {
	path string
}

func Config(path string) *ConfigOption {
	return &ConfigOption{path}
}

type DirOption struct {
	path string
}

func Dir(path string) *DirOption {
	return &DirOption{path}
}

type DirOrPlanOption struct {
	path string
}

func DirOrPlan(path string) *DirOrPlanOption {
	return &DirOrPlanOption{path}
}

// named to prevent conflict with DestroyOption interface
type DestroyFlagOption struct {
	destroy bool
}

func Destroy(destroy bool) *DestroyFlagOption {
	return &DestroyFlagOption{destroy}
}

type ForceCopyOption struct {
	forceCopy bool
}

func ForceCopy(forceCopy bool) *ForceCopyOption {
	return &ForceCopyOption{forceCopy}
}

type FromModuleOption struct {
	source string
}

func FromModule(source string) *FromModuleOption {
	return &FromModuleOption{source}
}

type GetOption struct {
	get bool
}

func Get(get bool) *GetOption {
	return &GetOption{get}
}

type GetPluginsOption struct {
	getPlugins bool
}

func GetPlugins(getPlugins bool) *GetPluginsOption {
	return &GetPluginsOption{getPlugins}
}

type LockOption struct {
	lock bool
}

func Lock(lock bool) *LockOption {
	return &LockOption{lock}
}

type LockTimeoutOption struct {
	timeout string
}

func LockTimeout(lockTimeout string) *LockTimeoutOption {
	return &LockTimeoutOption{lockTimeout}
}

type OutOption struct {
	path string
}

func Out(path string) *OutOption {
	return &OutOption{path}
}

type ParallelismOption struct {
	parallelism int
}

func Parallelism(n int) *ParallelismOption {
	return &ParallelismOption{n}
}

type PluginDirOption struct {
	pluginDir string
}

func PluginDir(pluginDir string) *PluginDirOption {
	return &PluginDirOption{pluginDir}
}

type ReconfigureOption struct {
	reconfigure bool
}

func Reconfigure(reconfigure bool) *ReconfigureOption {
	return &ReconfigureOption{reconfigure}
}

type RefreshOption struct {
	refresh bool
}

func Refresh(refresh bool) *RefreshOption {
	return &RefreshOption{refresh}
}

type StateOption struct {
	path string
}

func State(path string) *StateOption {
	return &StateOption{path}
}

type StateOutOption struct {
	path string
}

func StateOut(path string) *StateOutOption {
	return &StateOutOption{path}
}

type TargetOption struct {
	target string
}

func Target(resource string) *TargetOption {
	return &TargetOption{resource}
}

type UpgradeOption struct {
	upgrade bool
}

func Upgrade(upgrade bool) *UpgradeOption {
	return &UpgradeOption{upgrade}
}

type VarOption struct {
	assignment string
}

func Var(assignment string) *VarOption {
	return &VarOption{assignment}
}

type VarFileOption struct {
	path string
}

func VarFile(path string) *VarFileOption {
	return &VarFileOption{path}
}

type VerifyPluginsOption struct {
	verifyPlugins bool
}

func VerifyPlugins(verifyPlugins bool) *VerifyPluginsOption {
	return &VerifyPluginsOption{verifyPlugins}
}
