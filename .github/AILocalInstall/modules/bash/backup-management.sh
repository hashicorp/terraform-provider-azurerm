#!/bin/bash

#==============================================================================
# Backup Management Module
#==============================================================================
#
# This module handles all backup-related operations including creating,
# validating, and restoring backups of VS Code settings.
#

# Create timestamped backup
create_backup() {
    local settings_path="$1"
    local backup_dir
    backup_dir="$(get_backup_directory)"
    
    # Create backup directory if it doesn't exist
    mkdir -p "$backup_dir"
    
    # Check if settings file exists
    if [[ ! -f "$settings_path" ]]; then
        log_info "No existing settings file to backup"
        return 0
    fi
    
    # Validate JSON before backup
    if ! is_valid_json "$settings_path"; then
        log_warning "Existing settings file is not valid JSON, creating fake backup marker"
        create_fake_backup_marker "$backup_dir"
        return 0
    fi
    
    # Create timestamped backup
    local timestamp
    timestamp=$(date +"%Y%m%d_%H%M%S")
    local backup_file="$backup_dir/settings_backup_${timestamp}.json"
    
    if cp "$settings_path" "$backup_file"; then
        log_success "Created backup: $backup_file"
        
        # Create backup info file
        local info_file="$backup_dir/settings_backup_${timestamp}.info"
        cat > "$info_file" << INFO_EOF
{
    "timestamp": "$(date -Iseconds)",
    "original_path": "$settings_path",
    "backup_type": "complete",
    "original_size": $(wc -c < "$settings_path"),
    "validation": "passed"
}
INFO_EOF
        
        return 0
    else
        log_error "Failed to create backup"
        return 1
    fi
}

# Create fake backup marker for corrupted files
create_fake_backup_marker() {
    local backup_dir="$1"
    local timestamp
    timestamp=$(date +"%Y%m%d_%H%M%S")
    local marker_file="$backup_dir/settings_backup_${timestamp}.json"
    
    # Create minimal valid JSON as marker
    echo '{}' > "$marker_file"
    
    # Create info file indicating this is a fake backup
    local info_file="$backup_dir/settings_backup_${timestamp}.info"
    cat > "$info_file" << INFO_EOF
{
    "timestamp": "$(date -Iseconds)",
    "original_path": "corrupted",
    "backup_type": "fake",
    "original_size": 0,
    "validation": "failed"
}
INFO_EOF
    
    log_info "Created fake backup marker for corrupted settings"
}

# Get most recent backup file
get_most_recent_backup() {
    local backup_dir
    backup_dir="$(get_backup_directory)"
    
    if [[ ! -d "$backup_dir" ]]; then
        return 1
    fi
    
    # Find most recent backup file
    local latest_backup
    latest_backup=$(find "$backup_dir" -name "settings_backup_*.json" -type f -printf '%T@ %p\n' 2>/dev/null | sort -n | tail -1 | cut -d' ' -f2-)
    
    if [[ -n "$latest_backup" && -f "$latest_backup" ]]; then
        echo "$latest_backup"
        return 0
    else
        return 1
    fi
}

# Get backup info
get_backup_info() {
    local backup_file="$1"
    local info_file="${backup_file%.json}.info"
    
    if [[ -f "$info_file" ]]; then
        cat "$info_file"
    else
        # Create default info for backups without info files
        cat << INFO_EOF
{
    "timestamp": "unknown",
    "original_path": "unknown",
    "backup_type": "legacy",
    "original_size": $(wc -c < "$backup_file" 2>/dev/null || echo 0),
    "validation": "unknown"
}
INFO_EOF
    fi
}

# Restore from backup
restore_from_backup() {
    local settings_path="$1"
    local backup_file
    
    backup_file=$(get_most_recent_backup)
    if [[ -z "$backup_file" ]]; then
        log_error "No backup found to restore from"
        return 1
    fi
    
    # Get backup info
    local backup_info
    backup_info=$(get_backup_info "$backup_file")
    local backup_type
    backup_type=$(echo "$backup_info" | jq -r '.backup_type // "unknown"')
    
    case "$backup_type" in
        "complete")
            log_info "Restoring from complete backup: $backup_file"
            if cp "$backup_file" "$settings_path"; then
                log_success "Settings restored successfully"
                return 0
            else
                log_error "Failed to restore from backup"
                return 1
            fi
            ;;
        "fake")
            log_info "Backup was created for corrupted file, removing current settings"
            if [[ -f "$settings_path" ]]; then
                rm "$settings_path"
                log_success "Corrupted settings file removed"
            fi
            return 0
            ;;
        *)
            log_warning "Unknown backup type, attempting restore anyway"
            if cp "$backup_file" "$settings_path"; then
                log_success "Settings restored from unknown backup type"
                return 0
            else
                log_error "Failed to restore from backup"
                return 1
            fi
            ;;
    esac
}

# Clean backup directory
clean_backup_directory() {
    local backup_dir
    backup_dir="$(get_backup_directory)"
    
    if [[ -d "$backup_dir" ]]; then
        if rm -rf "$backup_dir"; then
            log_success "Cleaned up backup directory: $backup_dir"
            return 0
        else
            log_error "Failed to clean backup directory: $backup_dir"
            return 1
        fi
    else
        log_info "No backup directory to clean"
        return 0
    fi
}

# List available backups
list_backups() {
    local backup_dir
    backup_dir="$(get_backup_directory)"
    
    if [[ ! -d "$backup_dir" ]]; then
        echo "No backups found"
        return 1
    fi
    
    local backup_files
    backup_files=($(find "$backup_dir" -name "settings_backup_*.json" -type f | sort))
    
    if [[ ${#backup_files[@]} -eq 0 ]]; then
        echo "No backups found"
        return 1
    fi
    
    echo "Available backups:"
    for backup_file in "${backup_files[@]}"; do
        local info
        info=$(get_backup_info "$backup_file")
        local timestamp
        local backup_type
        local original_size
        
        timestamp=$(echo "$info" | jq -r '.timestamp // "unknown"')
        backup_type=$(echo "$info" | jq -r '.backup_type // "unknown"')
        original_size=$(echo "$info" | jq -r '.original_size // 0')
        
        printf "  %s (%s, %s, %d bytes)\n" \
            "$(basename "$backup_file")" \
            "$timestamp" \
            "$backup_type" \
            "$original_size"
    done
    
    return 0
}
