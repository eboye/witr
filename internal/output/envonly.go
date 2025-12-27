package output

import (
	"fmt"
	"strings"

	"github.com/pranshuparmar/witr/pkg/model"
)

// sensitivePatterns contains environment variable patterns that should be redacted
var sensitivePatterns = []string{
	"PASSWORD", "SECRET", "TOKEN", "KEY", "API_KEY", "PRIVATE_KEY",
	"AWS_SECRET", "AWS_ACCESS", "DATABASE_URL", "DB_PASSWORD", "CREDENTIAL",
	"AUTH", "PASSPHRASE", "CERTIFICATE", "SSL_CERT", "TLS_CERT",
}

// isSensitive checks if an environment variable might contain sensitive data
func isSensitive(envVar string) bool {
	// Split on = to get the variable name
	parts := strings.SplitN(envVar, "=", 2)
	if len(parts) < 1 {
		return false
	}
	
	varName := strings.ToUpper(parts[0])
	for _, pattern := range sensitivePatterns {
		if strings.Contains(varName, pattern) {
			return true
		}
	}
	return false
}

// redactEnvVar redacts the value of a sensitive environment variable
func redactEnvVar(envVar string) string {
	parts := strings.SplitN(envVar, "=", 2)
	if len(parts) < 2 {
		return envVar
	}
	return parts[0] + "=***REDACTED***"
}

// RenderEnvOnly prints only the command and environment variables for a process
func RenderEnvOnly(proc model.Process, colorEnabled bool) {
	RenderEnvOnlyWithRedaction(proc, colorEnabled, true)
}

// RenderEnvOnlyWithRedaction prints environment variables with optional redaction
func RenderEnvOnlyWithRedaction(proc model.Process, colorEnabled bool, redactSensitive bool) {
	colorResetEnv := ""
	colorBlueEnv := ""
	colorRedEnv := ""
	colorGreenEnv := ""
	colorYellowEnv := ""
	if colorEnabled {
		colorResetEnv = "\033[0m"
		colorBlueEnv = "\033[34m"
		colorRedEnv = "\033[31m"
		colorGreenEnv = "\033[32m"
		colorYellowEnv = "\033[33m"
	}
	fmt.Printf("%sCommand%s     : %s\n", colorGreenEnv, colorResetEnv, proc.Cmdline)
	if len(proc.Env) > 0 {
		fmt.Printf("%sEnvironment%s :\n", colorBlueEnv, colorResetEnv)
		sensitiveCount := 0
		for _, env := range proc.Env {
			displayEnv := env
			if redactSensitive && isSensitive(env) {
				displayEnv = redactEnvVar(env)
				sensitiveCount++
			}
			fmt.Printf("  %s\n", displayEnv)
		}
		if redactSensitive && sensitiveCount > 0 {
			fmt.Printf("\n%s[%d sensitive environment variable(s) redacted]%s\n", colorYellowEnv, sensitiveCount, colorResetEnv)
			fmt.Printf("%sUse --show-secrets to display all variables.%s\n", colorYellowEnv, colorResetEnv)
		}
	} else {
		fmt.Printf("%sNo environment variables found.%s\n", colorRedEnv, colorResetEnv)
	}
}
