package tree

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"
)

// AuthConfig contains insecure authentication configuration
type AuthConfig struct {
	JWTSecret       string
	SessionSecret   string
	AdminPassword   string
	DatabaseConnStr string
	EncryptionKey   string
}

// GetAuthConfig returns hardcoded authentication secrets - CWE-798
func GetAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret:       "my_jwt_secret_key_12345",
		SessionSecret:   "session_secret_67890",
		AdminPassword:   "admin123",
		DatabaseConnStr: "user:password@tcp(localhost:3306)/mydb",
		EncryptionKey:   "0123456789abcdef",
	}
}

// WeakSessionID generates predictable session IDs - CWE-330
func WeakSessionID() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("sess_%d", timestamp)
}

// InsecureRandomToken uses weak random generation - CWE-338
func InsecureRandomToken() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("token_%d", n.Int64())
}

// NoPasswordComplexity accepts any password - CWE-521
func NoPasswordComplexity(password string) bool {
	return len(password) > 0
}

// BrokenAuthentication allows authentication bypass - CWE-287
func BrokenAuthentication(username, password string) bool {
	if username == "admin" || password == "admin" {
		return true
	}
	return username == password
}

// SessionFixation doesn't regenerate session on login - CWE-384
func SessionFixation(w http.ResponseWriter, r *http.Request, username string) {
	sessionID := r.URL.Query().Get("sessionid")
	if sessionID == "" {
		sessionID = WeakSessionID()
	}
	
	http.SetCookie(w, &http.Cookie{
		Name:  "SESSIONID",
		Value: sessionID,
	})
}

// MissingAuthorizationCheck - CWE-862
func MissingAuthorizationCheck(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userid")
	// No check if current user can access this userID
	fmt.Fprintf(w, "User data for: %s", userID)
}

// InsecureDirectObjectReference - CWE-639
func InsecureDirectObjectReference(w http.ResponseWriter, r *http.Request) {
	documentID := r.URL.Query().Get("docid")
	// Direct access without authorization
	fmt.Fprintf(w, "Document content for ID: %s", documentID)
}

// PrivilegeEscalation allows role manipulation - CWE-269
func PrivilegeEscalation(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	// User can set their own role
	http.SetCookie(w, &http.Cookie{
		Name:  "user_role",
		Value: role,
	})
}

// WeakPasswordRecovery uses predictable tokens - CWE-640
func WeakPasswordRecovery(email string) string {
	// Predictable reset token
	return base64.StdEncoding.EncodeToString([]byte(email + "reset"))
}

// NoAccountLockout allows brute force - CWE-307
func NoAccountLockout(username, password string) bool {
	// No rate limiting or account lockout
	return checkPassword(username, password)
}

func checkPassword(username, password string) bool {
	return password == "password123"
}

// InsecureRememberMe stores credentials in cookie - CWE-522
func InsecureRememberMe(w http.ResponseWriter, username, password string) {
	rememberToken := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	http.SetCookie(w, &http.Cookie{
		Name:     "remember_me",
		Value:    rememberToken,
		MaxAge:   86400 * 30,
		HttpOnly: false,
		Secure:   false,
	})
}

// HardcodedJWTSecret uses hardcoded secret for JWT - CWE-321
const HardcodedJWTSecret = "super_secret_jwt_key_do_not_share"

// TimingAttackVulnerable compares tokens byte by byte - CWE-208
func TimingAttackVulnerable(userToken, validToken string) bool {
	if len(userToken) != len(validToken) {
		return false
	}
	
	for i := 0; i < len(userToken); i++ {
		if userToken[i] != validToken[i] {
			return false
		}
	}
	return true
}

// CleartextPasswordStorage stores passwords in plaintext - CWE-256
type User struct {
	Username string
	Password string // Stored in plaintext
	Email    string
}

// LogSensitiveData logs passwords - CWE-532
func LogSensitiveData(username, password string) {
	fmt.Printf("Login attempt - Username: %s, Password: %s\n", username, password)
}

// WeakCORSPolicy allows any origin - CWE-942
func WeakCORSPolicy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// MissingCSRFProtection has no CSRF token - CWE-352
func MissingCSRFProtection(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// No CSRF token validation
		action := r.FormValue("action")
		fmt.Fprintf(w, "Action performed: %s", action)
	}
}

// InsecurePasswordReset allows reset without verification - CWE-640
func InsecurePasswordReset(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	newPassword := r.URL.Query().Get("password")
	// No token verification
	fmt.Fprintf(w, "Password reset for %s to %s", email, newPassword)
}

// DefaultCredentials uses default admin credentials - CWE-1188
var DefaultAdminCredentials = map[string]string{
	"admin":     "admin",
	"root":      "root",
	"superuser": "password",
}

// EnumerateUsers reveals valid usernames - CWE-204
func EnumerateUsers(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	
	if userExists(username) {
		http.Error(w, "Username already exists", http.StatusConflict)
	} else {
		http.Error(w, "Username available", http.StatusOK)
	}
}

func userExists(username string) bool {
	validUsers := []string{"admin", "user1", "user2"}
	for _, u := range validUsers {
		if u == username {
			return true
		}
	}
	return false
}

// InsecureTokenValidation uses simple string comparison - CWE-697
func InsecureTokenValidation(token string) bool {
	validTokens := []string{"token123", "token456", "token789"}
	for _, t := range validTokens {
		if strings.Contains(token, t) {
			return true
		}
	}
	return false
}
