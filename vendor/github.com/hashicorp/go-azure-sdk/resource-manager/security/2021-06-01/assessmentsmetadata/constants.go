package assessmentsmetadata

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssessmentType string

const (
	AssessmentTypeBuiltIn         AssessmentType = "BuiltIn"
	AssessmentTypeCustomPolicy    AssessmentType = "CustomPolicy"
	AssessmentTypeCustomerManaged AssessmentType = "CustomerManaged"
	AssessmentTypeVerifiedPartner AssessmentType = "VerifiedPartner"
)

func PossibleValuesForAssessmentType() []string {
	return []string{
		string(AssessmentTypeBuiltIn),
		string(AssessmentTypeCustomPolicy),
		string(AssessmentTypeCustomerManaged),
		string(AssessmentTypeVerifiedPartner),
	}
}

func parseAssessmentType(input string) (*AssessmentType, error) {
	vals := map[string]AssessmentType{
		"builtin":         AssessmentTypeBuiltIn,
		"custompolicy":    AssessmentTypeCustomPolicy,
		"customermanaged": AssessmentTypeCustomerManaged,
		"verifiedpartner": AssessmentTypeVerifiedPartner,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssessmentType(input)
	return &out, nil
}

type Categories string

const (
	CategoriesCompute           Categories = "Compute"
	CategoriesData              Categories = "Data"
	CategoriesIdentityAndAccess Categories = "IdentityAndAccess"
	CategoriesIoT               Categories = "IoT"
	CategoriesNetworking        Categories = "Networking"
)

func PossibleValuesForCategories() []string {
	return []string{
		string(CategoriesCompute),
		string(CategoriesData),
		string(CategoriesIdentityAndAccess),
		string(CategoriesIoT),
		string(CategoriesNetworking),
	}
}

func parseCategories(input string) (*Categories, error) {
	vals := map[string]Categories{
		"compute":           CategoriesCompute,
		"data":              CategoriesData,
		"identityandaccess": CategoriesIdentityAndAccess,
		"iot":               CategoriesIoT,
		"networking":        CategoriesNetworking,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Categories(input)
	return &out, nil
}

type ImplementationEffort string

const (
	ImplementationEffortHigh     ImplementationEffort = "High"
	ImplementationEffortLow      ImplementationEffort = "Low"
	ImplementationEffortModerate ImplementationEffort = "Moderate"
)

func PossibleValuesForImplementationEffort() []string {
	return []string{
		string(ImplementationEffortHigh),
		string(ImplementationEffortLow),
		string(ImplementationEffortModerate),
	}
}

func parseImplementationEffort(input string) (*ImplementationEffort, error) {
	vals := map[string]ImplementationEffort{
		"high":     ImplementationEffortHigh,
		"low":      ImplementationEffortLow,
		"moderate": ImplementationEffortModerate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImplementationEffort(input)
	return &out, nil
}

type Severity string

const (
	SeverityHigh   Severity = "High"
	SeverityLow    Severity = "Low"
	SeverityMedium Severity = "Medium"
)

func PossibleValuesForSeverity() []string {
	return []string{
		string(SeverityHigh),
		string(SeverityLow),
		string(SeverityMedium),
	}
}

func parseSeverity(input string) (*Severity, error) {
	vals := map[string]Severity{
		"high":   SeverityHigh,
		"low":    SeverityLow,
		"medium": SeverityMedium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Severity(input)
	return &out, nil
}

type Tactics string

const (
	TacticsCollection          Tactics = "Collection"
	TacticsCommandAndControl   Tactics = "Command and Control"
	TacticsCredentialAccess    Tactics = "Credential Access"
	TacticsDefenseEvasion      Tactics = "Defense Evasion"
	TacticsDiscovery           Tactics = "Discovery"
	TacticsExecution           Tactics = "Execution"
	TacticsExfiltration        Tactics = "Exfiltration"
	TacticsImpact              Tactics = "Impact"
	TacticsInitialAccess       Tactics = "Initial Access"
	TacticsLateralMovement     Tactics = "Lateral Movement"
	TacticsPersistence         Tactics = "Persistence"
	TacticsPrivilegeEscalation Tactics = "Privilege Escalation"
	TacticsReconnaissance      Tactics = "Reconnaissance"
	TacticsResourceDevelopment Tactics = "Resource Development"
)

func PossibleValuesForTactics() []string {
	return []string{
		string(TacticsCollection),
		string(TacticsCommandAndControl),
		string(TacticsCredentialAccess),
		string(TacticsDefenseEvasion),
		string(TacticsDiscovery),
		string(TacticsExecution),
		string(TacticsExfiltration),
		string(TacticsImpact),
		string(TacticsInitialAccess),
		string(TacticsLateralMovement),
		string(TacticsPersistence),
		string(TacticsPrivilegeEscalation),
		string(TacticsReconnaissance),
		string(TacticsResourceDevelopment),
	}
}

func parseTactics(input string) (*Tactics, error) {
	vals := map[string]Tactics{
		"collection":           TacticsCollection,
		"command and control":  TacticsCommandAndControl,
		"credential access":    TacticsCredentialAccess,
		"defense evasion":      TacticsDefenseEvasion,
		"discovery":            TacticsDiscovery,
		"execution":            TacticsExecution,
		"exfiltration":         TacticsExfiltration,
		"impact":               TacticsImpact,
		"initial access":       TacticsInitialAccess,
		"lateral movement":     TacticsLateralMovement,
		"persistence":          TacticsPersistence,
		"privilege escalation": TacticsPrivilegeEscalation,
		"reconnaissance":       TacticsReconnaissance,
		"resource development": TacticsResourceDevelopment,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Tactics(input)
	return &out, nil
}

type Techniques string

const (
	TechniquesAbuseElevationControlMechanism          Techniques = "Abuse Elevation Control Mechanism"
	TechniquesAccessTokenManipulation                 Techniques = "Access Token Manipulation"
	TechniquesAccountDiscovery                        Techniques = "Account Discovery"
	TechniquesAccountManipulation                     Techniques = "Account Manipulation"
	TechniquesActiveScanning                          Techniques = "Active Scanning"
	TechniquesApplicationLayerProtocol                Techniques = "Application Layer Protocol"
	TechniquesAudioCapture                            Techniques = "Audio Capture"
	TechniquesBootOrLogonAutostartExecution           Techniques = "Boot or Logon Autostart Execution"
	TechniquesBootOrLogonInitializationScripts        Techniques = "Boot or Logon Initialization Scripts"
	TechniquesBruteForce                              Techniques = "Brute Force"
	TechniquesCloudInfrastructureDiscovery            Techniques = "Cloud Infrastructure Discovery"
	TechniquesCloudServiceDashboard                   Techniques = "Cloud Service Dashboard"
	TechniquesCloudServiceDiscovery                   Techniques = "Cloud Service Discovery"
	TechniquesCommandAndScriptingInterpreter          Techniques = "Command and Scripting Interpreter"
	TechniquesCompromiseClientSoftwareBinary          Techniques = "Compromise Client Software Binary"
	TechniquesCompromiseInfrastructure                Techniques = "Compromise Infrastructure"
	TechniquesContainerAndResourceDiscovery           Techniques = "Container and Resource Discovery"
	TechniquesCreateAccount                           Techniques = "Create Account"
	TechniquesCreateOrModifySystemProcess             Techniques = "Create or Modify System Process"
	TechniquesCredentialsFromPasswordStores           Techniques = "Credentials from Password Stores"
	TechniquesDataDestruction                         Techniques = "Data Destruction"
	TechniquesDataEncryptedForImpact                  Techniques = "Data Encrypted for Impact"
	TechniquesDataFromCloudStorageObject              Techniques = "Data from Cloud Storage Object"
	TechniquesDataFromConfigurationRepository         Techniques = "Data from Configuration Repository"
	TechniquesDataFromInformationRepositories         Techniques = "Data from Information Repositories"
	TechniquesDataFromLocalSystem                     Techniques = "Data from Local System"
	TechniquesDataManipulation                        Techniques = "Data Manipulation"
	TechniquesDataStaged                              Techniques = "Data Staged"
	TechniquesDefacement                              Techniques = "Defacement"
	TechniquesDeobfuscateDecodeFilesOrInformation     Techniques = "Deobfuscate/Decode Files or Information"
	TechniquesDiskWipe                                Techniques = "Disk Wipe"
	TechniquesDomainTrustDiscovery                    Techniques = "Domain Trust Discovery"
	TechniquesDriveNegativebyCompromise               Techniques = "Drive-by Compromise"
	TechniquesDynamicResolution                       Techniques = "Dynamic Resolution"
	TechniquesEndpointDenialOfService                 Techniques = "Endpoint Denial of Service"
	TechniquesEventTriggeredExecution                 Techniques = "Event Triggered Execution"
	TechniquesExfiltrationOverAlternativeProtocol     Techniques = "Exfiltration Over Alternative Protocol"
	TechniquesExploitPublicNegativeFacingApplication  Techniques = "Exploit Public-Facing Application"
	TechniquesExploitationForClientExecution          Techniques = "Exploitation for Client Execution"
	TechniquesExploitationForCredentialAccess         Techniques = "Exploitation for Credential Access"
	TechniquesExploitationForDefenseEvasion           Techniques = "Exploitation for Defense Evasion"
	TechniquesExploitationForPrivilegeEscalation      Techniques = "Exploitation for Privilege Escalation"
	TechniquesExploitationOfRemoteServices            Techniques = "Exploitation of Remote Services"
	TechniquesExternalRemoteServices                  Techniques = "External Remote Services"
	TechniquesFallbackChannels                        Techniques = "Fallback Channels"
	TechniquesFileAndDirectoryDiscovery               Techniques = "File and Directory Discovery"
	TechniquesFileAndDirectoryPermissionsModification Techniques = "File and Directory Permissions Modification"
	TechniquesGatherVictimNetworkInformation          Techniques = "Gather Victim Network Information"
	TechniquesHideArtifacts                           Techniques = "Hide Artifacts"
	TechniquesHijackExecutionFlow                     Techniques = "Hijack Execution Flow"
	TechniquesImpairDefenses                          Techniques = "Impair Defenses"
	TechniquesImplantContainerImage                   Techniques = "Implant Container Image"
	TechniquesIndicatorRemovalOnHost                  Techniques = "Indicator Removal on Host"
	TechniquesIndirectCommandExecution                Techniques = "Indirect Command Execution"
	TechniquesIngressToolTransfer                     Techniques = "Ingress Tool Transfer"
	TechniquesInputCapture                            Techniques = "Input Capture"
	TechniquesInterNegativeProcessCommunication       Techniques = "Inter-Process Communication"
	TechniquesLateralToolTransfer                     Techniques = "Lateral Tool Transfer"
	TechniquesManNegativeinNegativetheNegativeMiddle  Techniques = "Man-in-the-Middle"
	TechniquesMasquerading                            Techniques = "Masquerading"
	TechniquesModifyAuthenticationProcess             Techniques = "Modify Authentication Process"
	TechniquesModifyRegistry                          Techniques = "Modify Registry"
	TechniquesNetworkDenialOfService                  Techniques = "Network Denial of Service"
	TechniquesNetworkServiceScanning                  Techniques = "Network Service Scanning"
	TechniquesNetworkSniffing                         Techniques = "Network Sniffing"
	TechniquesNonNegativeApplicationLayerProtocol     Techniques = "Non-Application Layer Protocol"
	TechniquesNonNegativeStandardPort                 Techniques = "Non-Standard Port"
	TechniquesOSCredentialDumping                     Techniques = "OS Credential Dumping"
	TechniquesObfuscatedFilesOrInformation            Techniques = "Obfuscated Files or Information"
	TechniquesObtainCapabilities                      Techniques = "Obtain Capabilities"
	TechniquesOfficeApplicationStartup                Techniques = "Office Application Startup"
	TechniquesPermissionGroupsDiscovery               Techniques = "Permission Groups Discovery"
	TechniquesPhishing                                Techniques = "Phishing"
	TechniquesPreNegativeOSBoot                       Techniques = "Pre-OS Boot"
	TechniquesProcessDiscovery                        Techniques = "Process Discovery"
	TechniquesProcessInjection                        Techniques = "Process Injection"
	TechniquesProtocolTunneling                       Techniques = "Protocol Tunneling"
	TechniquesProxy                                   Techniques = "Proxy"
	TechniquesQueryRegistry                           Techniques = "Query Registry"
	TechniquesRemoteAccessSoftware                    Techniques = "Remote Access Software"
	TechniquesRemoteServiceSessionHijacking           Techniques = "Remote Service Session Hijacking"
	TechniquesRemoteServices                          Techniques = "Remote Services"
	TechniquesRemoteSystemDiscovery                   Techniques = "Remote System Discovery"
	TechniquesResourceHijacking                       Techniques = "Resource Hijacking"
	TechniquesSQLStoredProcedures                     Techniques = "SQL Stored Procedures"
	TechniquesScheduledTaskJob                        Techniques = "Scheduled Task/Job"
	TechniquesScreenCapture                           Techniques = "Screen Capture"
	TechniquesSearchVictimNegativeOwnedWebsites       Techniques = "Search Victim-Owned Websites"
	TechniquesServerSoftwareComponent                 Techniques = "Server Software Component"
	TechniquesServiceStop                             Techniques = "Service Stop"
	TechniquesSignedBinaryProxyExecution              Techniques = "Signed Binary Proxy Execution"
	TechniquesSoftwareDeploymentTools                 Techniques = "Software Deployment Tools"
	TechniquesStealOrForgeKerberosTickets             Techniques = "Steal or Forge Kerberos Tickets"
	TechniquesSubvertTrustControls                    Techniques = "Subvert Trust Controls"
	TechniquesSupplyChainCompromise                   Techniques = "Supply Chain Compromise"
	TechniquesSystemInformationDiscovery              Techniques = "System Information Discovery"
	TechniquesTaintSharedContent                      Techniques = "Taint Shared Content"
	TechniquesTrafficSignaling                        Techniques = "Traffic Signaling"
	TechniquesTransferDataToCloudAccount              Techniques = "Transfer Data to Cloud Account"
	TechniquesTrustedRelationship                     Techniques = "Trusted Relationship"
	TechniquesUnsecuredCredentials                    Techniques = "Unsecured Credentials"
	TechniquesUserExecution                           Techniques = "User Execution"
	TechniquesValidAccounts                           Techniques = "Valid Accounts"
	TechniquesWindowsManagementInstrumentation        Techniques = "Windows Management Instrumentation"
)

func PossibleValuesForTechniques() []string {
	return []string{
		string(TechniquesAbuseElevationControlMechanism),
		string(TechniquesAccessTokenManipulation),
		string(TechniquesAccountDiscovery),
		string(TechniquesAccountManipulation),
		string(TechniquesActiveScanning),
		string(TechniquesApplicationLayerProtocol),
		string(TechniquesAudioCapture),
		string(TechniquesBootOrLogonAutostartExecution),
		string(TechniquesBootOrLogonInitializationScripts),
		string(TechniquesBruteForce),
		string(TechniquesCloudInfrastructureDiscovery),
		string(TechniquesCloudServiceDashboard),
		string(TechniquesCloudServiceDiscovery),
		string(TechniquesCommandAndScriptingInterpreter),
		string(TechniquesCompromiseClientSoftwareBinary),
		string(TechniquesCompromiseInfrastructure),
		string(TechniquesContainerAndResourceDiscovery),
		string(TechniquesCreateAccount),
		string(TechniquesCreateOrModifySystemProcess),
		string(TechniquesCredentialsFromPasswordStores),
		string(TechniquesDataDestruction),
		string(TechniquesDataEncryptedForImpact),
		string(TechniquesDataFromCloudStorageObject),
		string(TechniquesDataFromConfigurationRepository),
		string(TechniquesDataFromInformationRepositories),
		string(TechniquesDataFromLocalSystem),
		string(TechniquesDataManipulation),
		string(TechniquesDataStaged),
		string(TechniquesDefacement),
		string(TechniquesDeobfuscateDecodeFilesOrInformation),
		string(TechniquesDiskWipe),
		string(TechniquesDomainTrustDiscovery),
		string(TechniquesDriveNegativebyCompromise),
		string(TechniquesDynamicResolution),
		string(TechniquesEndpointDenialOfService),
		string(TechniquesEventTriggeredExecution),
		string(TechniquesExfiltrationOverAlternativeProtocol),
		string(TechniquesExploitPublicNegativeFacingApplication),
		string(TechniquesExploitationForClientExecution),
		string(TechniquesExploitationForCredentialAccess),
		string(TechniquesExploitationForDefenseEvasion),
		string(TechniquesExploitationForPrivilegeEscalation),
		string(TechniquesExploitationOfRemoteServices),
		string(TechniquesExternalRemoteServices),
		string(TechniquesFallbackChannels),
		string(TechniquesFileAndDirectoryDiscovery),
		string(TechniquesFileAndDirectoryPermissionsModification),
		string(TechniquesGatherVictimNetworkInformation),
		string(TechniquesHideArtifacts),
		string(TechniquesHijackExecutionFlow),
		string(TechniquesImpairDefenses),
		string(TechniquesImplantContainerImage),
		string(TechniquesIndicatorRemovalOnHost),
		string(TechniquesIndirectCommandExecution),
		string(TechniquesIngressToolTransfer),
		string(TechniquesInputCapture),
		string(TechniquesInterNegativeProcessCommunication),
		string(TechniquesLateralToolTransfer),
		string(TechniquesManNegativeinNegativetheNegativeMiddle),
		string(TechniquesMasquerading),
		string(TechniquesModifyAuthenticationProcess),
		string(TechniquesModifyRegistry),
		string(TechniquesNetworkDenialOfService),
		string(TechniquesNetworkServiceScanning),
		string(TechniquesNetworkSniffing),
		string(TechniquesNonNegativeApplicationLayerProtocol),
		string(TechniquesNonNegativeStandardPort),
		string(TechniquesOSCredentialDumping),
		string(TechniquesObfuscatedFilesOrInformation),
		string(TechniquesObtainCapabilities),
		string(TechniquesOfficeApplicationStartup),
		string(TechniquesPermissionGroupsDiscovery),
		string(TechniquesPhishing),
		string(TechniquesPreNegativeOSBoot),
		string(TechniquesProcessDiscovery),
		string(TechniquesProcessInjection),
		string(TechniquesProtocolTunneling),
		string(TechniquesProxy),
		string(TechniquesQueryRegistry),
		string(TechniquesRemoteAccessSoftware),
		string(TechniquesRemoteServiceSessionHijacking),
		string(TechniquesRemoteServices),
		string(TechniquesRemoteSystemDiscovery),
		string(TechniquesResourceHijacking),
		string(TechniquesSQLStoredProcedures),
		string(TechniquesScheduledTaskJob),
		string(TechniquesScreenCapture),
		string(TechniquesSearchVictimNegativeOwnedWebsites),
		string(TechniquesServerSoftwareComponent),
		string(TechniquesServiceStop),
		string(TechniquesSignedBinaryProxyExecution),
		string(TechniquesSoftwareDeploymentTools),
		string(TechniquesStealOrForgeKerberosTickets),
		string(TechniquesSubvertTrustControls),
		string(TechniquesSupplyChainCompromise),
		string(TechniquesSystemInformationDiscovery),
		string(TechniquesTaintSharedContent),
		string(TechniquesTrafficSignaling),
		string(TechniquesTransferDataToCloudAccount),
		string(TechniquesTrustedRelationship),
		string(TechniquesUnsecuredCredentials),
		string(TechniquesUserExecution),
		string(TechniquesValidAccounts),
		string(TechniquesWindowsManagementInstrumentation),
	}
}

func parseTechniques(input string) (*Techniques, error) {
	vals := map[string]Techniques{
		"abuse elevation control mechanism":           TechniquesAbuseElevationControlMechanism,
		"access token manipulation":                   TechniquesAccessTokenManipulation,
		"account discovery":                           TechniquesAccountDiscovery,
		"account manipulation":                        TechniquesAccountManipulation,
		"active scanning":                             TechniquesActiveScanning,
		"application layer protocol":                  TechniquesApplicationLayerProtocol,
		"audio capture":                               TechniquesAudioCapture,
		"boot or logon autostart execution":           TechniquesBootOrLogonAutostartExecution,
		"boot or logon initialization scripts":        TechniquesBootOrLogonInitializationScripts,
		"brute force":                                 TechniquesBruteForce,
		"cloud infrastructure discovery":              TechniquesCloudInfrastructureDiscovery,
		"cloud service dashboard":                     TechniquesCloudServiceDashboard,
		"cloud service discovery":                     TechniquesCloudServiceDiscovery,
		"command and scripting interpreter":           TechniquesCommandAndScriptingInterpreter,
		"compromise client software binary":           TechniquesCompromiseClientSoftwareBinary,
		"compromise infrastructure":                   TechniquesCompromiseInfrastructure,
		"container and resource discovery":            TechniquesContainerAndResourceDiscovery,
		"create account":                              TechniquesCreateAccount,
		"create or modify system process":             TechniquesCreateOrModifySystemProcess,
		"credentials from password stores":            TechniquesCredentialsFromPasswordStores,
		"data destruction":                            TechniquesDataDestruction,
		"data encrypted for impact":                   TechniquesDataEncryptedForImpact,
		"data from cloud storage object":              TechniquesDataFromCloudStorageObject,
		"data from configuration repository":          TechniquesDataFromConfigurationRepository,
		"data from information repositories":          TechniquesDataFromInformationRepositories,
		"data from local system":                      TechniquesDataFromLocalSystem,
		"data manipulation":                           TechniquesDataManipulation,
		"data staged":                                 TechniquesDataStaged,
		"defacement":                                  TechniquesDefacement,
		"deobfuscate/decode files or information":     TechniquesDeobfuscateDecodeFilesOrInformation,
		"disk wipe":                                   TechniquesDiskWipe,
		"domain trust discovery":                      TechniquesDomainTrustDiscovery,
		"drive-by compromise":                         TechniquesDriveNegativebyCompromise,
		"dynamic resolution":                          TechniquesDynamicResolution,
		"endpoint denial of service":                  TechniquesEndpointDenialOfService,
		"event triggered execution":                   TechniquesEventTriggeredExecution,
		"exfiltration over alternative protocol":      TechniquesExfiltrationOverAlternativeProtocol,
		"exploit public-facing application":           TechniquesExploitPublicNegativeFacingApplication,
		"exploitation for client execution":           TechniquesExploitationForClientExecution,
		"exploitation for credential access":          TechniquesExploitationForCredentialAccess,
		"exploitation for defense evasion":            TechniquesExploitationForDefenseEvasion,
		"exploitation for privilege escalation":       TechniquesExploitationForPrivilegeEscalation,
		"exploitation of remote services":             TechniquesExploitationOfRemoteServices,
		"external remote services":                    TechniquesExternalRemoteServices,
		"fallback channels":                           TechniquesFallbackChannels,
		"file and directory discovery":                TechniquesFileAndDirectoryDiscovery,
		"file and directory permissions modification": TechniquesFileAndDirectoryPermissionsModification,
		"gather victim network information":           TechniquesGatherVictimNetworkInformation,
		"hide artifacts":                              TechniquesHideArtifacts,
		"hijack execution flow":                       TechniquesHijackExecutionFlow,
		"impair defenses":                             TechniquesImpairDefenses,
		"implant container image":                     TechniquesImplantContainerImage,
		"indicator removal on host":                   TechniquesIndicatorRemovalOnHost,
		"indirect command execution":                  TechniquesIndirectCommandExecution,
		"ingress tool transfer":                       TechniquesIngressToolTransfer,
		"input capture":                               TechniquesInputCapture,
		"inter-process communication":                 TechniquesInterNegativeProcessCommunication,
		"lateral tool transfer":                       TechniquesLateralToolTransfer,
		"man-in-the-middle":                           TechniquesManNegativeinNegativetheNegativeMiddle,
		"masquerading":                                TechniquesMasquerading,
		"modify authentication process":               TechniquesModifyAuthenticationProcess,
		"modify registry":                             TechniquesModifyRegistry,
		"network denial of service":                   TechniquesNetworkDenialOfService,
		"network service scanning":                    TechniquesNetworkServiceScanning,
		"network sniffing":                            TechniquesNetworkSniffing,
		"non-application layer protocol":              TechniquesNonNegativeApplicationLayerProtocol,
		"non-standard port":                           TechniquesNonNegativeStandardPort,
		"os credential dumping":                       TechniquesOSCredentialDumping,
		"obfuscated files or information":             TechniquesObfuscatedFilesOrInformation,
		"obtain capabilities":                         TechniquesObtainCapabilities,
		"office application startup":                  TechniquesOfficeApplicationStartup,
		"permission groups discovery":                 TechniquesPermissionGroupsDiscovery,
		"phishing":                                    TechniquesPhishing,
		"pre-os boot":                                 TechniquesPreNegativeOSBoot,
		"process discovery":                           TechniquesProcessDiscovery,
		"process injection":                           TechniquesProcessInjection,
		"protocol tunneling":                          TechniquesProtocolTunneling,
		"proxy":                                       TechniquesProxy,
		"query registry":                              TechniquesQueryRegistry,
		"remote access software":                      TechniquesRemoteAccessSoftware,
		"remote service session hijacking":            TechniquesRemoteServiceSessionHijacking,
		"remote services":                             TechniquesRemoteServices,
		"remote system discovery":                     TechniquesRemoteSystemDiscovery,
		"resource hijacking":                          TechniquesResourceHijacking,
		"sql stored procedures":                       TechniquesSQLStoredProcedures,
		"scheduled task/job":                          TechniquesScheduledTaskJob,
		"screen capture":                              TechniquesScreenCapture,
		"search victim-owned websites":                TechniquesSearchVictimNegativeOwnedWebsites,
		"server software component":                   TechniquesServerSoftwareComponent,
		"service stop":                                TechniquesServiceStop,
		"signed binary proxy execution":               TechniquesSignedBinaryProxyExecution,
		"software deployment tools":                   TechniquesSoftwareDeploymentTools,
		"steal or forge kerberos tickets":             TechniquesStealOrForgeKerberosTickets,
		"subvert trust controls":                      TechniquesSubvertTrustControls,
		"supply chain compromise":                     TechniquesSupplyChainCompromise,
		"system information discovery":                TechniquesSystemInformationDiscovery,
		"taint shared content":                        TechniquesTaintSharedContent,
		"traffic signaling":                           TechniquesTrafficSignaling,
		"transfer data to cloud account":              TechniquesTransferDataToCloudAccount,
		"trusted relationship":                        TechniquesTrustedRelationship,
		"unsecured credentials":                       TechniquesUnsecuredCredentials,
		"user execution":                              TechniquesUserExecution,
		"valid accounts":                              TechniquesValidAccounts,
		"windows management instrumentation":          TechniquesWindowsManagementInstrumentation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Techniques(input)
	return &out, nil
}

type Threats string

const (
	ThreatsAccountBreach        Threats = "accountBreach"
	ThreatsDataExfiltration     Threats = "dataExfiltration"
	ThreatsDataSpillage         Threats = "dataSpillage"
	ThreatsDenialOfService      Threats = "denialOfService"
	ThreatsElevationOfPrivilege Threats = "elevationOfPrivilege"
	ThreatsMaliciousInsider     Threats = "maliciousInsider"
	ThreatsMissingCoverage      Threats = "missingCoverage"
	ThreatsThreatResistance     Threats = "threatResistance"
)

func PossibleValuesForThreats() []string {
	return []string{
		string(ThreatsAccountBreach),
		string(ThreatsDataExfiltration),
		string(ThreatsDataSpillage),
		string(ThreatsDenialOfService),
		string(ThreatsElevationOfPrivilege),
		string(ThreatsMaliciousInsider),
		string(ThreatsMissingCoverage),
		string(ThreatsThreatResistance),
	}
}

func parseThreats(input string) (*Threats, error) {
	vals := map[string]Threats{
		"accountbreach":        ThreatsAccountBreach,
		"dataexfiltration":     ThreatsDataExfiltration,
		"dataspillage":         ThreatsDataSpillage,
		"denialofservice":      ThreatsDenialOfService,
		"elevationofprivilege": ThreatsElevationOfPrivilege,
		"maliciousinsider":     ThreatsMaliciousInsider,
		"missingcoverage":      ThreatsMissingCoverage,
		"threatresistance":     ThreatsThreatResistance,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Threats(input)
	return &out, nil
}

type UserImpact string

const (
	UserImpactHigh     UserImpact = "High"
	UserImpactLow      UserImpact = "Low"
	UserImpactModerate UserImpact = "Moderate"
)

func PossibleValuesForUserImpact() []string {
	return []string{
		string(UserImpactHigh),
		string(UserImpactLow),
		string(UserImpactModerate),
	}
}

func parseUserImpact(input string) (*UserImpact, error) {
	vals := map[string]UserImpact{
		"high":     UserImpactHigh,
		"low":      UserImpactLow,
		"moderate": UserImpactModerate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UserImpact(input)
	return &out, nil
}
