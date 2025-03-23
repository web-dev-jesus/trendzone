package util

import (
	"net/url"
	"regexp"
	"strings"
)

// SanitizeURL removes sensitive information like API keys from URLs for logging
func SanitizeURL(rawURL string) string {
	// Parse URL
	u, err := url.Parse(rawURL)
	if err != nil {
		// If can't parse, do best effort string replacement
		return redactAPIKey(rawURL)
	}

	// Get query parameters
	q := u.Query()

	// Redact sensitive query parameters
	sensitiveParams := []string{"key", "apikey", "api_key", "password", "token", "access_token"}
	for _, param := range sensitiveParams {
		if q.Get(param) != "" {
			q.Set(param, "REDACTED")
		}
	}

	// Set sanitized query parameters
	u.RawQuery = q.Encode()

	// Return sanitized URL
	return u.String()
}

// redactAPIKey performs a string-based redaction of API keys for when URL parsing fails
func redactAPIKey(input string) string {
	// Define patterns for common API key parameters
	patterns := []string{
		`key=[^&]+`,
		`apikey=[^&]+`,
		`api_key=[^&]+`,
		`token=[^&]+`,
		`access_token=[^&]+`,
	}

	result := input
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		parts := strings.SplitN(pattern, "=", 2)
		if len(parts) == 2 {
			param := parts[0]
			result = re.ReplaceAllString(result, param+"=REDACTED")
		}
	}

	return result
}

// ValidateInput validates user input against potential security issues
func ValidateInput(input string) bool {
	// Check for SQL injection patterns
	sqlInjectionPatterns := []string{
		`(?i)'\s*OR\s*'1'='1`,
		`(?i)--`,
		`(?i);\s*DROP\s+TABLE`,
		`(?i);\s*DELETE\s+FROM`,
		`(?i)UNION\s+SELECT`,
	}

	for _, pattern := range sqlInjectionPatterns {
		match, _ := regexp.MatchString(pattern, input)
		if match {
			return false
		}
	}

	// Check for NoSQL injection patterns
	noSqlInjectionPatterns := []string{
		`(?i)\$where`,
		`(?i)\$ne`,
		`(?i)\$gt`,
		`(?i)\$lt`,
		`(?i)\$exists`,
	}

	for _, pattern := range noSqlInjectionPatterns {
		match, _ := regexp.MatchString(pattern, input)
		if match {
			return false
		}
	}

	// Check for command injection
	commandInjectionPatterns := []string{
		`;`,
		`&&`,
		`\|\|`,
		`\$\(`,
		"`",
	}

	for _, pattern := range commandInjectionPatterns {
		match, _ := regexp.MatchString(pattern, input)
		if match {
			return false
		}
	}

	return true
}

// HashPassword returns a secure hash of a password
// This is a placeholder - in a real implementation, use bcrypt or similar
func HashPassword(password string) (string, error) {
	// In a real implementation:
	// return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return "hashed_password", nil
}

// GenerateSecureToken generates a secure random token
// This is a placeholder - in a real implementation, use crypto/rand
func GenerateSecureToken(length int) (string, error) {
	// In a real implementation:
	// Use crypto/rand to generate secure random bytes
	// Convert to a string representation
	return "secure_random_token", nil
}
