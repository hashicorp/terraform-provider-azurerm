// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_app_attach_package

type VirtualDesktopAppAttachPackageModel struct {
	Name                       string                        `tfschema:"name"`
	ResourceGroupName          string                        `tfschema:"resource_group_name"`
	Location                   string                        `tfschema:"location"`
	DisplayName                string                        `tfschema:"display_name"`
	HostPoolIds                []string                      `tfschema:"host_pool_ids"`
	MsixPackageName            string                        `tfschema:"msix_package_name"`
	StorageShareFileId         string                        `tfschema:"storage_share_file_id"`
	HealthCheckStatusOnFailure string                        `tfschema:"health_check_status_on_failure"`
	RegisterAtLogOnEnabled     bool                          `tfschema:"register_at_log_on_enabled"`
	StateEnabled               bool                          `tfschema:"state_enabled"`
	Tags                       map[string]string             `tfschema:"tags"`
	LastUpdated                string                        `tfschema:"last_updated"`
	PackageApplications        []MsixPackageApplicationModel `tfschema:"package_applications"`
	PackageFamilyName          string                        `tfschema:"package_family_name"`
	PackageName                string                        `tfschema:"package_name"`
	PackageRelativePath        string                        `tfschema:"package_relative_path"`
	Version                    string                        `tfschema:"version"`
}

type MsixPackageApplicationModel struct {
	AppId          string `tfschema:"app_id"`
	AppUserModelID string `tfschema:"app_user_model_id"`
	Description    string `tfschema:"description"`
	FriendlyName   string `tfschema:"friendly_name"`
	IconImageName  string `tfschema:"icon_image_name"`
	RawIcon        string `tfschema:"raw_icon"`
	RawPng         string `tfschema:"raw_png"`
}

func (r Resource) ModelObject() any {
	return &VirtualDesktopAppAttachPackageModel{}
}
