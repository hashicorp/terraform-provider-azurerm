---
applyTo: "internal/**/*.go"
description: Security and compliance patterns for the Terraform AzureRM provider including input validation, credential management, and security best practices.
---

# Security & Compliance Patterns

Security and compliance patterns for the Terraform AzureRM provider including input validation, credential management, and security best practices.

**Quick navigation:** [🔐 Input Validation](#🔐-input-validation-and-sanitization) | [🔑 Credential Management](#🔑-credential-management) | [🛡️ Security Patterns](#🛡️-security-patterns) | [📋 Compliance](#📋-compliance-requirements)

## 🔐 Input Validation and Sanitization

### Preventing Injection Attacks

```go
func ValidateSecureResourceName(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)

    // Prevent injection attacks
    dangerousPatterns := []string{"'", "\"", ";", "--", "/*", "*/", "<", ">"}
    for _, pattern := range dangerousPatterns {
        if strings.Contains(value, pattern) {
            errors = append(errors, fmt.Errorf("property `%s` cannot contain potentially unsafe characters", k))
            return warnings, errors
        }
    }

    // Validate length constraints
    if len(value) < 1 || len(value) > 64 {
        errors = append(errors, fmt.Errorf("property `%s` must be between 1 and 64 characters", k))
        return warnings, errors
    }

    // Azure-specific character validation
    allowedPattern := regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
    if !allowedPattern.MatchString(value) {
        errors = append(errors, fmt.Errorf("property `%s` can only contain alphanumeric characters, hyphens, and underscores", k))
        return warnings, errors
    }

    return warnings, errors
}
```

### SQL Injection Prevention

```go
func ValidateSQLResourceName(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)

    // Check for SQL injection patterns
    sqlInjectionPatterns := []string{
        "'", "\"", ";", "--", "/*", "*/", "xp_", "sp_", "exec", "execute",
        "select", "insert", "update", "delete", "drop", "create", "alter",
    }

    lowerValue := strings.ToLower(value)
    for _, pattern := range sqlInjectionPatterns {
        if strings.Contains(lowerValue, pattern) {
            errors = append(errors, fmt.Errorf("property `%s` cannot contain potentially unsafe characters or SQL keywords", k))
            return warnings, errors
        }
    }

    return warnings, errors
}
```

### Path Traversal Prevention

```go
func ValidateFilePath(v interface{}, k string) (warnings []string, errors []error) {
    value := v.(string)

    // Prevent path traversal attacks
    if strings.Contains(value, "..") {
        errors = append(errors, fmt.Errorf("property `%s` cannot contain path traversal sequences", k))
        return warnings, errors
    }

    // Prevent absolute paths in certain contexts
    if strings.HasPrefix(value, "/") || strings.Contains(value, ":") {
        errors = append(errors, fmt.Errorf("property `%s` must be a relative path", k))
        return warnings, errors
    }

    return warnings, errors
}
```

## 🔑 Credential Management

### Secure Environment Variable Handling

```go
func validateTestCredentials() error {
    requiredVars := []string{
        "ARM_SUBSCRIPTION_ID",
        "ARM_CLIENT_ID",
        "ARM_CLIENT_SECRET",
        "ARM_TENANT_ID",
    }

    for _, envVar := range requiredVars {
        if value := os.Getenv(envVar); value == "" {
            return fmt.Errorf("required environment variable %s is not set", envVar)
        }
    }
    return nil
}

// Secure credential validation
func validateCredentialFormat(credential string, credType string) error {
    switch credType {
    case "subscription_id", "tenant_id":
        // UUID format validation
        uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
        if !uuidPattern.MatchString(credential) {
            return fmt.Errorf("invalid %s format", credType)
        }
    case "client_secret":
        // Minimum entropy requirements
        if len(credential) < 32 {
            return fmt.Errorf("client secret too short")
        }
    }
    return nil
}
```

### Sensitive Data Handling

```go
// Never log sensitive information
func logSecurely(message string, sensitiveData map[string]interface{}) {
    // Redact sensitive fields
    safe := make(map[string]interface{})
    for k, v := range sensitiveData {
        if isSensitiveField(k) {
            safe[k] = "[REDACTED]"
        } else {
            safe[k] = v
        }
    }

    log.Printf("%s: %+v", message, safe)
}

func isSensitiveField(fieldName string) bool {
    sensitiveFields := []string{
        "password", "secret", "key", "token", "credential",
        "connection_string", "sas_token", "access_key",
    }

    lower := strings.ToLower(fieldName)
    for _, sensitive := range sensitiveFields {
        if strings.Contains(lower, sensitive) {
            return true
        }
    }
    return false
}
```

## 🛡️ Security Patterns

### Secure Resource Creation

```go
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // Validate security requirements
            if err := validateSecurityRequirements(metadata); err != nil {
                return fmt.Errorf("security validation failed: %+v", err)
            }

            // Apply security defaults
            properties := applySecurityDefaults(resourceProperties)

            // Audit logging
            auditLog := map[string]interface{}{
                "operation":    "create",
                "resource_id":  id.String(),
                "timestamp":    time.Now().UTC(),
                "user_context": getUserContext(metadata),
            }
            logSecurely("Resource creation", auditLog)

            return nil
        },
    }
}

func validateSecurityRequirements(metadata sdk.ResourceMetaData) error {
    // Implement security policy validation
    if !isAuthorizedOperation(metadata) {
        return fmt.Errorf("operation not authorized")
    }

    if !meetsComplianceRequirements(metadata) {
        return fmt.Errorf("operation does not meet compliance requirements")
    }

    return nil
}
```

### Encryption and Data Protection

```go
func applyEncryptionDefaults(properties *azureapi.ResourceProperties) {
    // Always enable encryption at rest
    if properties.Encryption == nil {
        properties.Encryption = &azureapi.EncryptionSettings{
            Enabled:                 pointer.To(true),
            EncryptionAtRestEnabled: pointer.To(true),
            KeySource:              pointer.To(azureapi.KeySourceMicrosoftStorage),
        }
    }

    // Enforce TLS minimum version
    if properties.TLSSettings == nil {
        properties.TLSSettings = &azureapi.TLSSettings{
            MinimumTLSVersion: pointer.To(azureapi.TLSVersion12),
        }
    }
}

func validateEncryptionSettings(settings *azureapi.EncryptionSettings) error {
    if settings.Enabled != nil && !*settings.Enabled {
        return fmt.Errorf("encryption cannot be disabled for compliance reasons")
    }

    if settings.KeySource != nil && *settings.KeySource == azureapi.KeySourceNone {
        return fmt.Errorf("encryption key source must be specified")
    }

    return nil
}
```

### Network Security

```go
func validateNetworkSecurity(networkConfig *azureapi.NetworkConfiguration) error {
    // Ensure HTTPS endpoints
    if networkConfig.PublicEndpoint != nil {
        endpoint := *networkConfig.PublicEndpoint
        if !strings.HasPrefix(endpoint, "https://") {
            return fmt.Errorf("public endpoints must use HTTPS")
        }
    }

    // Validate firewall rules
    if networkConfig.FirewallRules != nil {
        for _, rule := range *networkConfig.FirewallRules {
            if err := validateFirewallRule(rule); err != nil {
                return fmt.Errorf("invalid firewall rule: %+v", err)
            }
        }
    }

    return nil
}

func validateFirewallRule(rule azureapi.FirewallRule) error {
    // Prevent overly permissive rules
    if rule.SourceIPRange != nil && *rule.SourceIPRange == "0.0.0.0/0" {
        return fmt.Errorf("cannot allow access from all IP addresses")
    }

    // Ensure secure ports
    if rule.Port != nil {
        port := *rule.Port
        if port < 1024 && port != 443 && port != 80 {
            return fmt.Errorf("system ports below 1024 are restricted")
        }
    }

    return nil
}
```

## 📋 Compliance Requirements

### Audit Logging

```go
type AuditEvent struct {
    Timestamp    time.Time              `json:"timestamp"`
    Operation    string                 `json:"operation"`
    ResourceType string                 `json:"resource_type"`
    ResourceID   string                 `json:"resource_id"`
    UserContext  map[string]interface{} `json:"user_context"`
    Changes      map[string]interface{} `json:"changes"`
    Result       string                 `json:"result"`
}

func logAuditEvent(event AuditEvent) {
    // Format for compliance systems
    auditJSON, _ := json.Marshal(event)

    // Send to audit system (implementation depends on requirements)
    sendToAuditSystem(auditJSON)

    // Local logging for debugging
    log.Printf("[AUDIT] %s", auditJSON)
}

func trackResourceChanges(old, new interface{}) map[string]interface{} {
    changes := make(map[string]interface{})

    // Calculate diff between old and new state
    // Implementation depends on specific requirements

    return changes
}
```

### Data Residency and Sovereignty

```go
func validateDataResidency(location string, dataClassification string) error {
    // Define data residency requirements
    restrictions := map[string][]string{
        "sensitive": {"US", "EU"},
        "public":    {"*"},
    }

    allowedRegions, exists := restrictions[dataClassification]
    if !exists {
        return fmt.Errorf("unknown data classification: %s", dataClassification)
    }

    if allowedRegions[0] != "*" {
        region := extractRegionFromLocation(location)
        if !contains(allowedRegions, region) {
            return fmt.Errorf("data classification %s not allowed in region %s", dataClassification, region)
        }
    }

    return nil
}

func extractRegionFromLocation(location string) string {
    // Map Azure locations to compliance regions
    regionMap := map[string]string{
        "eastus":     "US",
        "westus":     "US",
        "westeurope": "EU",
        "northeurope": "EU",
    }

    if region, exists := regionMap[strings.ToLower(location)]; exists {
        return region
    }

    return "Unknown"
}
```

### Compliance Validation

```go
func validateComplianceRequirements(ctx context.Context, resource interface{}) error {
    // GDPR compliance
    if err := validateGDPRCompliance(resource); err != nil {
        return fmt.Errorf("GDPR compliance failed: %+v", err)
    }

    // SOC 2 compliance
    if err := validateSOC2Compliance(resource); err != nil {
        return fmt.Errorf("SOC 2 compliance failed: %+v", err)
    }

    // HIPAA compliance (if applicable)
    if isHealthcareContext(ctx) {
        if err := validateHIPAACompliance(resource); err != nil {
            return fmt.Errorf("HIPAA compliance failed: %+v", err)
        }
    }

    return nil
}

func validateGDPRCompliance(resource interface{}) error {
    // Ensure data protection by design
    if !hasDataProtectionByDesign(resource) {
        return fmt.Errorf("resource must implement data protection by design")
    }

    // Validate data retention policies
    if !hasDataRetentionPolicy(resource) {
        return fmt.Errorf("resource must have data retention policy")
    }

    return nil
}
```
---
[⬆️ Back to top](#security--compliance-patterns)
