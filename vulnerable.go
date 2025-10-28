package tree

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// VulnerableHashPassword uses weak MD5 hashing for passwords
func VulnerableHashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// VulnerableExecuteCommand executes shell commands without sanitization
func VulnerableExecuteCommand(userInput string) error {
	cmd := exec.Command("sh", "-c", userInput)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

// VulnerableFileRead reads files without path validation
func VulnerableFileRead(filename string) ([]byte, error) {
	// Path traversal vulnerability - no validation
	return os.ReadFile(filename)
}

// VulnerableHTTPHandler handles HTTP requests with SQL injection vulnerability
func VulnerableHTTPHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	
	// SQL injection vulnerability - direct string concatenation
	query := "SELECT * FROM users WHERE username = '" + username + "'"
	
	fmt.Fprintf(w, "Executing query: %s", query)
}

// VulnerableFileUpload saves uploaded files without validation
func VulnerableFileUpload(filename string, content io.Reader) error {
	// No file type validation or size limits
	dst, err := os.Create(filepath.Join("/tmp", filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, content)
	return err
}

// VulnerableAPIKey exposes hardcoded credentials
const (
	APIKey    = "sk-1234567890abcdef"
	DBPassword = "admin123"
	SecretToken = "my-secret-token-12345"
)

// VulnerableRandomNumber generates predictable random numbers
func VulnerableRandomNumber() int {
	// Using weak random number generation
	return 42 // Completely predictable
}

// VulnerableDeserialize deserializes data without validation
func VulnerableDeserialize(data []byte) error {
	// Unsafe deserialization - no type checking
	// This would be vulnerable if using encoding/gob or similar
	return nil
}

// VulnerableXSS returns user input without sanitization
func VulnerableXSS(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("message")
	// XSS vulnerability - no HTML escaping
	fmt.Fprintf(w, "<html><body>%s</body></html>", userInput)
}

// VulnerableRaceCondition has a race condition vulnerability
var globalCounter int

func VulnerableRaceCondition() {
	// Race condition - no synchronization
	globalCounter++
}
