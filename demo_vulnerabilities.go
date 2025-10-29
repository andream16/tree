package tree

import (
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

// Hardcoded credentials - CWE-798
const (
	DatabasePassword = "P@ssw0rd123!"
	APISecret        = "secret_key_1234567890abcdefghijklmnop"
	PrivateKey       = "-----BEGIN PRIVATE KEY-----\nEXAMPLE_KEY_DATA_HERE..."
	AWSAccessKeyID   = "EXAMPLE_ACCESS_KEY_ID_12345"
	AWSSecretKey     = "example_secret_access_key_67890"
)

// WeakHashPassword uses MD5 for password hashing - CWE-327
func WeakHashPassword(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

// InsecureSHA1Hash uses deprecated SHA1 - CWE-327
func InsecureSHA1Hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

// SQLInjectionVulnerable - CWE-89
func SQLInjectionVulnerable(db *sql.DB, username string) (*sql.Rows, error) {
	query := "SELECT * FROM users WHERE username = '" + username + "'"
	return db.Query(query)
}

// CommandInjection - CWE-78
func CommandInjection(filename string) (string, error) {
	cmd := exec.Command("sh", "-c", "cat "+filename)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// PathTraversal - CWE-22
func PathTraversal(userPath string) ([]byte, error) {
	fullPath := "/var/www/uploads/" + userPath
	return os.ReadFile(fullPath)
}

// UnsafeFileUpload - CWE-434
func UnsafeFileUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("upload")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create("/tmp/" + header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)
	fmt.Fprintf(w, "File uploaded: %s", header.Filename)
}

// XSSVulnerable - CWE-79
func XSSVulnerable(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "<html><body><h1>Hello, %s!</h1></body></html>", name)
}

// OpenRedirect - CWE-601
func OpenRedirect(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("redirect")
	http.Redirect(w, r, url, http.StatusFound)
}

// SSRFVulnerable - CWE-918
func SSRFVulnerable(targetURL string) ([]byte, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// InsecureCookie - CWE-614, CWE-1004
func InsecureCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: false,
		Secure:   false,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

// RaceCondition - CWE-362
var GlobalBalance int

func RaceCondition(amount int) {
	GlobalBalance = GlobalBalance + amount
}

// NullPointerDereference - CWE-476
func NullPointerDereference(data *string) string {
	return *data
}

// IntegerOverflow - CWE-190
func IntegerOverflow(a, b int32) int32 {
	return a + b
}

// UnhandledError - CWE-391
func UnhandledError(filename string) []byte {
	data, _ := os.ReadFile(filename)
	return data
}

// MemoryLeak - CWE-401
func MemoryLeak() {
	data := make([]byte, 1024*1024*100) // 100MB
	_ = data
}

// InfiniteLoop - CWE-835
func InfiniteLoop(input string) {
	for input != "" {
		// Missing break condition
	}
}
