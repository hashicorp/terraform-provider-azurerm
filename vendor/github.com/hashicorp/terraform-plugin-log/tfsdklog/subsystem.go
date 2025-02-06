package tfsdklog

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/internal/hclogutils"
	"github.com/hashicorp/terraform-plugin-log/internal/logging"
)

// NewSubsystem returns a new context.Context that contains a subsystem logger
// configured with the passed options, named after the subsystem argument.
//
// Subsystem loggers allow different areas of a plugin codebase to use
// different logging levels, giving developers more fine-grained control over
// what is logging and with what verbosity. They're best utilized for logical
// concerns that are sometimes helpful to log, but may generate unwanted noise
// at other times.
//
// The only Options supported for subsystems are the Options for setting the
// level and additional location offset of the logger.
func NewSubsystem(ctx context.Context, subsystem string, options ...logging.Option) context.Context {
	logger := logging.GetSDKRootLogger(ctx)

	if logger == nil {
		// this essentially should never happen in production
		// the root logger for provider code should be injected
		// by whatever  the provider developer is using, so
		// really this is only likely in unit tests, at most
		// so just making this a no-op is fine
		return ctx
	}

	rootLoggerOptions := logging.GetSDKRootLoggerOptions(ctx)
	subLoggerTFLoggerOpts := logging.ApplyLoggerOpts(options...)

	// If root logger options are not available,
	// fallback to creating a logger named like the given subsystem.
	// This will preserve the root logger options,
	// but cannot make changes beyond setting the level
	// due to limitations with the hclog.Logger interface.
	var subLogger hclog.Logger
	if rootLoggerOptions == nil {
		subLogger = logger.Named(subsystem)

		if subLoggerTFLoggerOpts.AdditionalLocationOffset != 1 {
			logger.Warn("Unable to create logging subsystem with AdditionalLocationOffset due to missing root logger options")
		}
	} else {
		subLoggerOptions := hclogutils.LoggerOptionsCopy(rootLoggerOptions)
		subLoggerOptions.Name = subLoggerOptions.Name + "." + subsystem

		if subLoggerTFLoggerOpts.AdditionalLocationOffset != 1 {
			subLoggerOptions.AdditionalLocationOffset = subLoggerTFLoggerOpts.AdditionalLocationOffset
		}

		subLogger = hclog.New(subLoggerOptions)
	}

	// Cache subsystem logger level outside context for performance reasons.
	subsystemLevelsMutex.Lock()

	subsystemLevels[subsystem] = subLoggerTFLoggerOpts.Level

	subsystemLevelsMutex.Unlock()

	// Set the configured log level
	if subLoggerTFLoggerOpts.Level != hclog.NoLevel {
		subLogger.SetLevel(subLoggerTFLoggerOpts.Level)
	}

	// Propagate root fields to the subsystem logger
	if subLoggerTFLoggerOpts.IncludeRootFields {
		loggerTFOpts := logging.GetSDKRootTFLoggerOpts(ctx)
		subLoggerTFLoggerOpts = logging.WithFields(loggerTFOpts.Fields)(subLoggerTFLoggerOpts)
	}

	// Set the subsystem LoggerOpts in the context
	ctx = logging.SetSDKSubsystemTFLoggerOpts(ctx, subsystem, subLoggerTFLoggerOpts)

	return logging.SetSDKSubsystemLogger(ctx, subsystem, subLogger)
}

// SubsystemSetField returns a new context.Context that has a modified logger for
// the specified subsystem in it which will include key and value as fields
// in all its log output.
//
// In case of the same key is used multiple times (i.e. key collision),
// the last one set is the one that gets persisted and then outputted with the logs.
func SubsystemSetField(ctx context.Context, subsystem, key string, value interface{}) context.Context {
	lOpts := logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem)

	// Copy to prevent slice/map aliasing issues.
	// Reference: https://github.com/hashicorp/terraform-plugin-log/issues/131
	lOpts = logging.WithField(key, value)(lOpts.Copy())

	return logging.SetSDKSubsystemTFLoggerOpts(ctx, subsystem, lOpts)
}

// SubsystemTrace logs `msg` at the trace level to the subsystem logger
// specified in `ctx`, with optional `additionalFields` structured key-value
// fields in the log output. Fields are shallow merged with any defined on the
// subsystem logger, e.g. by the `SubsystemSetField()` function, and across
// multiple maps.
func SubsystemTrace(ctx context.Context, subsystem, msg string, additionalFields ...map[string]interface{}) {
	if !subsystemWouldLog(subsystem, hclog.Trace) {
		return
	}

	logger := logging.GetSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		if logging.GetSDKRootLogger(ctx) == nil {
			// logging isn't set up, nothing we can do, just silently fail
			// this should basically never happen in production
			return
		}
		// create a new logger if one doesn't exist
		logger = logging.GetSDKSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem).With("new_logger_warning", logging.NewSDKSubsystemLoggerWarning)
	}

	additionalArgs, shouldOmit := logging.OmitOrMask(logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem), &msg, additionalFields)
	if shouldOmit {
		return
	}

	logger.Trace(msg, additionalArgs...)
}

// SubsystemDebug logs `msg` at the debug level to the subsystem logger
// specified in `ctx`, with optional `additionalFields` structured key-value
// fields in the log output. Fields are shallow merged with any defined on the
// subsystem logger, e.g. by the `SubsystemSetField()` function, and across
// multiple maps.
func SubsystemDebug(ctx context.Context, subsystem, msg string, additionalFields ...map[string]interface{}) {
	if !subsystemWouldLog(subsystem, hclog.Debug) {
		return
	}

	logger := logging.GetSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		if logging.GetSDKRootLogger(ctx) == nil {
			// logging isn't set up, nothing we can do, just silently fail
			// this should basically never happen in production
			return
		}
		// create a new logger if one doesn't exist
		logger = logging.GetSDKSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem).With("new_logger_warning", logging.NewSDKSubsystemLoggerWarning)
	}

	additionalArgs, shouldOmit := logging.OmitOrMask(logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem), &msg, additionalFields)
	if shouldOmit {
		return
	}

	logger.Debug(msg, additionalArgs...)
}

// SubsystemInfo logs `msg` at the info level to the subsystem logger
// specified in `ctx`, with optional `additionalFields` structured key-value
// fields in the log output. Fields are shallow merged with any defined on the
// subsystem logger, e.g. by the `SubsystemSetField()` function, and across
// multiple maps.
func SubsystemInfo(ctx context.Context, subsystem, msg string, additionalFields ...map[string]interface{}) {
	if !subsystemWouldLog(subsystem, hclog.Info) {
		return
	}

	logger := logging.GetSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		if logging.GetSDKRootLogger(ctx) == nil {
			// logging isn't set up, nothing we can do, just silently fail
			// this should basically never happen in production
			return
		}
		// create a new logger if one doesn't exist
		logger = logging.GetSDKSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem).With("new_logger_warning", logging.NewSDKSubsystemLoggerWarning)
	}

	additionalArgs, shouldOmit := logging.OmitOrMask(logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem), &msg, additionalFields)
	if shouldOmit {
		return
	}

	logger.Info(msg, additionalArgs...)
}

// SubsystemWarn logs `msg` at the warn level to the subsystem logger
// specified in `ctx`, with optional `additionalFields` structured key-value
// fields in the log output. Fields are shallow merged with any defined on the
// subsystem logger, e.g. by the `SubsystemSetField()` function, and across
// multiple maps.
func SubsystemWarn(ctx context.Context, subsystem, msg string, additionalFields ...map[string]interface{}) {
	if !subsystemWouldLog(subsystem, hclog.Warn) {
		return
	}

	logger := logging.GetSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		if logging.GetSDKRootLogger(ctx) == nil {
			// logging isn't set up, nothing we can do, just silently fail
			// this should basically never happen in production
			return
		}
		// create a new logger if one doesn't exist
		logger = logging.GetSDKSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem).With("new_logger_warning", logging.NewSDKSubsystemLoggerWarning)
	}

	additionalArgs, shouldOmit := logging.OmitOrMask(logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem), &msg, additionalFields)
	if shouldOmit {
		return
	}

	logger.Warn(msg, additionalArgs...)
}

// SubsystemError logs `msg` at the error level to the subsystem logger
// specified in `ctx`, with optional `additionalFields` structured key-value
// fields in the log output. Fields are shallow merged with any defined on the
// subsystem logger, e.g. by the `SubsystemSetField()` function, and across
// multiple maps.
func SubsystemError(ctx context.Context, subsystem, msg string, additionalFields ...map[string]interface{}) {
	if !subsystemWouldLog(subsystem, hclog.Error) {
		return
	}

	logger := logging.GetSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		if logging.GetSDKRootLogger(ctx) == nil {
			// logging isn't set up, nothing we can do, just silently fail
			// this should basically never happen in production
			return
		}
		// create a new logger if one doesn't exist
		logger = logging.GetSDKSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem).With("new_logger_warning", logging.NewSDKSubsystemLoggerWarning)
	}

	additionalArgs, shouldOmit := logging.OmitOrMask(logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem), &msg, additionalFields)
	if shouldOmit {
		return
	}

	logger.Error(msg, additionalArgs...)
}

// SubsystemOmitLogWithFieldKeys returns a new context.Context that has a modified logger
// that will omit to write any log when any of the given keys is found
// within its fields.
//
// Each call to this function is additive:
// the keys to omit by are added to the existing configuration.
//
// Example:
//
//	configuration = `['foo', 'baz']`
//
//	log1 = `{ msg = "...", fields = { 'foo': '...', 'bar': '...' }`  -> omitted
//	log2 = `{ msg = "...", fields = { 'bar': '...' }`                -> printed
//	log3 = `{ msg = "...", fields = { 'baz': '...', 'boo': '...' }`  -> omitted
func SubsystemOmitLogWithFieldKeys(ctx context.Context, subsystem string, keys ...string) context.Context {
	lOpts := logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem)

	// Copy to prevent slice/map aliasing issues.
	// Reference: https://github.com/hashicorp/terraform-plugin-log/issues/131
	lOpts = logging.WithOmitLogWithFieldKeys(keys...)(lOpts.Copy())

	return logging.SetSDKSubsystemTFLoggerOpts(ctx, subsystem, lOpts)
}

// SubsystemOmitLogWithMessageRegexes returns a new context.Context that has a modified logger
// that will omit to write any log that has a message matching any of the
// given *regexp.Regexp.
//
// Each call to this function is additive:
// the regexp to omit by are added to the existing configuration.
//
// Example:
//
//	configuration = `[regexp.MustCompile("(foo|bar)")]`
//
//	log1 = `{ msg = "banana apple foo", fields = {...}`     -> omitted
//	log2 = `{ msg = "pineapple mango", fields = {...}`      -> printed
//	log3 = `{ msg = "pineapple mango bar", fields = {...}`  -> omitted
func SubsystemOmitLogWithMessageRegexes(ctx context.Context, subsystem string, expressions ...*regexp.Regexp) context.Context {
	lOpts := logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem)

	// Copy to prevent slice/map aliasing issues.
	// Reference: https://github.com/hashicorp/terraform-plugin-log/issues/131
	lOpts = logging.WithOmitLogWithMessageRegexes(expressions...)(lOpts.Copy())

	return logging.SetSDKSubsystemTFLoggerOpts(ctx, subsystem, lOpts)
}

// SubsystemOmitLogWithMessageStrings  returns a new context.Context that has a modified logger
// that will omit to write any log that matches any of the given string.
//
// Each call to this function is additive:
// the string to omit by are added to the existing configuration.
//
// Example:
//
//	configuration = `['foo', 'bar']`
//
//	log1 = `{ msg = "banana apple foo", fields = {...}`     -> omitted
//	log2 = `{ msg = "pineapple mango", fields = {...}`      -> printed
//	log3 = `{ msg = "pineapple mango bar", fields = {...}`  -> omitted
func SubsystemOmitLogWithMessageStrings(ctx context.Context, subsystem string, matchingStrings ...string) context.Context {
	lOpts := logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem)

	// Copy to prevent slice/map aliasing issues.
	// Reference: https://github.com/hashicorp/terraform-plugin-log/issues/131
	lOpts = logging.WithOmitLogWithMessageStrings(matchingStrings...)(lOpts.Copy())

	return logging.SetSDKSubsystemTFLoggerOpts(ctx, subsystem, lOpts)
}

// SubsystemMaskFieldValuesWithFieldKeys returns a new context.Context that has a modified logger
// that masks (replaces) with asterisks (`***`) any field value where the
// key matches one of the given keys.
//
// Each call to this function is additive:
// the keys to mask by are added to the existing configuration.
//
// Example:
//
//	configuration = `['foo', 'baz']`
//
//	log1 = `{ msg = "...", fields = { 'foo': '***', 'bar': '...' }`  -> masked value
//	log2 = `{ msg = "...", fields = { 'bar': '...' }`                -> as-is value
//	log3 = `{ msg = "...", fields = { 'baz': '***', 'boo': '...' }`  -> masked value
func SubsystemMaskFieldValuesWithFieldKeys(ctx context.Context, subsystem string, keys ...string) context.Context {
	lOpts := logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem)

	// Copy to prevent slice/map aliasing issues.
	// Reference: https://github.com/hashicorp/terraform-plugin-log/issues/131
	lOpts = logging.WithMaskFieldValuesWithFieldKeys(keys...)(lOpts.Copy())

	return logging.SetSDKSubsystemTFLoggerOpts(ctx, subsystem, lOpts)
}

// SubsystemMaskAllFieldValuesRegexes returns a new context.Context that has a modified logger
// that masks (replaces) with asterisks (`***`) all field value substrings,
// matching one of the given *regexp.Regexp.
//
// Each call to this function is additive:
// the regexp to mask by are added to the existing configuration.
//
// Example:
//
//	configuration = `[regexp.MustCompile("(foo|bar)")]`
//
//	log1 = `{ msg = "...", fields = { 'k1': '***', 'k2': '***', 'k3': 'baz' }`  -> masked value
//	log2 = `{ msg = "...", fields = { 'k1': 'boo', 'k2': 'far', 'k3': 'baz' }`  -> as-is value
//	log2 = `{ msg = "...", fields = { 'k1': '*** *** baz' }`                    -> masked value
func SubsystemMaskAllFieldValuesRegexes(ctx context.Context, subsystem string, expressions ...*regexp.Regexp) context.Context {
	lOpts := logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem)

	// Copy to prevent slice/map aliasing issues.
	// Reference: https://github.com/hashicorp/terraform-plugin-log/issues/131
	lOpts = logging.WithMaskAllFieldValuesRegexes(expressions...)(lOpts.Copy())

	return logging.SetSDKSubsystemTFLoggerOpts(ctx, subsystem, lOpts)
}

// SubsystemMaskAllFieldValuesStrings returns a new context.Context that has a modified logger
// that masks (replaces) with asterisks (`***`) all field value substrings,
// equal to one of the given strings.
//
// Each call to this function is additive:
// the regexp to mask by are added to the existing configuration.
//
// Example:
//
//	configuration = `[regexp.MustCompile("(foo|bar)")]`
//
//	log1 = `{ msg = "...", fields = { 'k1': '***', 'k2': '***', 'k3': 'baz' }`  -> masked value
//	log2 = `{ msg = "...", fields = { 'k1': 'boo', 'k2': 'far', 'k3': 'baz' }`  -> as-is value
//	log2 = `{ msg = "...", fields = { 'k1': '*** *** baz' }`                    -> masked value
func SubsystemMaskAllFieldValuesStrings(ctx context.Context, subsystem string, matchingStrings ...string) context.Context {
	lOpts := logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem)

	// Copy to prevent slice/map aliasing issues.
	// Reference: https://github.com/hashicorp/terraform-plugin-log/issues/131
	lOpts = logging.WithMaskAllFieldValuesStrings(matchingStrings...)(lOpts.Copy())

	return logging.SetSDKSubsystemTFLoggerOpts(ctx, subsystem, lOpts)
}

// SubsystemMaskMessageRegexes returns a new context.Context that has a modified logger
// that masks (replaces) with asterisks (`***`) all message substrings,
// matching one of the given *regexp.Regexp.
//
// Each call to this function is additive:
// the regexp to mask by are added to the existing configuration.
//
// Example:
//
//	configuration = `[regexp.MustCompile("(foo|bar)")]`
//
//	log1 = `{ msg = "banana apple ***", fields = {...}`     -> masked portion
//	log2 = `{ msg = "pineapple mango", fields = {...}`      -> as-is
//	log3 = `{ msg = "pineapple mango ***", fields = {...}`  -> masked portion
func SubsystemMaskMessageRegexes(ctx context.Context, subsystem string, expressions ...*regexp.Regexp) context.Context {
	lOpts := logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem)

	// Copy to prevent slice/map aliasing issues.
	// Reference: https://github.com/hashicorp/terraform-plugin-log/issues/131
	lOpts = logging.WithMaskMessageRegexes(expressions...)(lOpts.Copy())

	return logging.SetSDKSubsystemTFLoggerOpts(ctx, subsystem, lOpts)
}

// SubsystemMaskMessageStrings returns a new context.Context that has a modified logger
// that masks (replace) with asterisks (`***`) all message substrings,
// equal to one of the given strings.
//
// Each call to this function is additive:
// the string to mask by are added to the existing configuration.
//
// Example:
//
//	configuration = `['foo', 'bar']`
//
//	log1 = `{ msg = "banana apple ***", fields = { 'k1': 'foo, bar, baz' }`  -> masked portion
//	log2 = `{ msg = "pineapple mango", fields = {...}`                       -> as-is
//	log3 = `{ msg = "pineapple mango ***", fields = {...}`                   -> masked portion
func SubsystemMaskMessageStrings(ctx context.Context, subsystem string, matchingStrings ...string) context.Context {
	lOpts := logging.GetSDKSubsystemTFLoggerOpts(ctx, subsystem)

	// Copy to prevent slice/map aliasing issues.
	// Reference: https://github.com/hashicorp/terraform-plugin-log/issues/131
	lOpts = logging.WithMaskMessageStrings(matchingStrings...)(lOpts.Copy())

	return logging.SetSDKSubsystemTFLoggerOpts(ctx, subsystem, lOpts)
}

// SubsystemMaskLogRegexes is a shortcut to invoke SubsystemMaskMessageRegexes and SubsystemMaskAllFieldValuesRegexes using the same input.
// Refer to those functions for details.
func SubsystemMaskLogRegexes(ctx context.Context, subsystem string, expressions ...*regexp.Regexp) context.Context {
	return SubsystemMaskMessageRegexes(SubsystemMaskAllFieldValuesRegexes(ctx, subsystem, expressions...), subsystem, expressions...)
}

// SubsystemMaskLogStrings is a shortcut to invoke SubsystemMaskMessageStrings and SubsystemMaskAllFieldValuesStrings using the same input.
// Refer to those functions for details.
func SubsystemMaskLogStrings(ctx context.Context, subsystem string, matchingStrings ...string) context.Context {
	return SubsystemMaskMessageStrings(SubsystemMaskAllFieldValuesStrings(ctx, subsystem, matchingStrings...), subsystem, matchingStrings...)
}
