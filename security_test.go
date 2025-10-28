package tree

import (
	"crypto/des"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

// InsecureConfig contains hardcoded sensitive information
type InsecureConfig struct {
	DatabaseURL  string
	AdminAPIKey  string
	JWTSecret    string
	AWSAccessKey string
	AWSSecretKey string
}

// GetInsecureConfig returns configuration with hardcoded secrets
func GetInsecureConfig() *InsecureConfig {
	return &InsecureConfig{
		DatabaseURL:  "postgres://admin:password123@localhost:5432/mydb",
		AdminAPIKey:  "AKIAIOSFODNN7EXAMPLE",
		JWTSecret:    "super-secret-jwt-key-12345",
		AWSAccessKey: "AKIAIOSFODNN7EXAMPLE",
		AWSSecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}
}

// VulnerableDesEncryption uses deprecated DES encryption
func VulnerableDesEncryption(key []byte, plaintext []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	ciphertext := make([]byte, len(plaintext))
	block.Encrypt(ciphertext, plaintext)
	return ciphertext, nil
}

// VulnerableRegexDOS is vulnerable to ReDoS attacks
func VulnerableRegexDOS(input string) bool {
	// Catastrophic backtracking pattern
	pattern := `^(a+)+$`
	matched, _ := regexp.MatchString(pattern, input)
	return matched
}

// VulnerableXMLParsing is vulnerable to XXE attacks
func VulnerableXMLParsing(xmlData []byte) error {
	type User struct {
		Name  string `xml:"name"`
		Email string `xml:"email"`
	}
	
	var user User
	// No protection against XXE
	return xml.Unmarshal(xmlData, &user)
}

// VulnerableLogInjection logs user input without sanitization
func VulnerableLogInjection(username string) {
	// Log injection vulnerability
	log.Printf("User login attempt: %s", username)
}

// VulnerableWeakRandom uses weak random number generation for security
func VulnerableWeakRandom() string {
	rand.Seed(time.Now().UnixNano())
	token := make([]byte, 16)
	for i := range token {
		token[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", token)
}

// VulnerableOpenRedirect allows open redirect attacks
func VulnerableOpenRedirect(w http.ResponseWriter, r *http.Request) {
	redirectURL := r.URL.Query().Get("redirect")
	// No validation of redirect URL
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// VulnerableSSRF is vulnerable to Server-Side Request Forgery
func VulnerableSSRF(url string) ([]byte, error) {
	// No URL validation - can access internal resources
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	return io.ReadAll(resp.Body)
}

// VulnerableDirectoryListing exposes directory contents
func VulnerableDirectoryListing(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	
	files, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	for _, file := range files {
		fmt.Fprintf(w, "%s\n", file.Name())
	}
}

// VulnerableNoSQLInjection is vulnerable to NoSQL injection
func VulnerableNoSQLInjection(username string) string {
	// MongoDB query with string concatenation
	query := fmt.Sprintf(`{"username": "%s"}`, username)
	return query
}

// VulnerableBufferOverflow attempts to write beyond buffer
func VulnerableBufferOverflow(input []byte) {
	buffer := make([]byte, 10)
	// No bounds checking
	copy(buffer, input)
}

// VulnerableTimingAttack compares secrets with timing vulnerability
func VulnerableTimingAttack(userToken, validToken string) bool {
	// Timing attack vulnerability - early return
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

// VulnerableUnvalidatedRedirect redirects without validation
func VulnerableUnvalidatedRedirect(w http.ResponseWriter, r *http.Request) {
	target := r.FormValue("url")
	w.Header().Set("Location", target)
	w.WriteHeader(http.StatusMovedPermanently)
}

// VulnerableLDAPInjection is vulnerable to LDAP injection
func VulnerableLDAPInjection(username string) string {
	// LDAP query with unsanitized input
	filter := fmt.Sprintf("(uid=%s)", username)
	return filter
}

// VulnerableInsecureDeserialization deserializes untrusted data
func VulnerableInsecureDeserialization(data []byte) (interface{}, error) {
	// Unsafe deserialization without type validation
	var result interface{}
	// This would be vulnerable with gob or other serialization
	return result, nil
}

// VulnerableResourceExhaustion has no rate limiting
func VulnerableResourceExhaustion(w http.ResponseWriter, r *http.Request) {
	// No rate limiting - vulnerable to DoS
	count := r.URL.Query().Get("count")
	
	for i := 0; i < 1000000; i++ {
		fmt.Fprintf(w, "Processing %s: %d\n", count, i)
	}
}

// VulnerableSQLQuery executes SQL with string formatting
func VulnerableSQLQuery(db *sql.DB, userID string) error {
	// SQL injection via string formatting
	query := fmt.Sprintf("DELETE FROM users WHERE id = %s", userID)
	_, err := db.Exec(query)
	return err
}

// VulnerableSessionFixation doesn't regenerate session IDs
func VulnerableSessionFixation(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	// Session fixation - accepts session ID from user
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	})
}

// VulnerableInsecureCookie sets cookies without security flags
func VulnerableInsecureCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "sensitive-token-value",
		HttpOnly: false, // Vulnerable to XSS
		Secure:   false, // Sent over HTTP
		SameSite: http.SameSiteNoneMode,
	})
}

// VulnerableMemoryLeak creates memory leaks
func VulnerableMemoryLeak() {
	// Memory leak - goroutine never exits
	go func() {
		data := make([]byte, 1024*1024) // 1MB
		for {
			_ = data
			time.Sleep(time.Second)
		}
	}()
}
