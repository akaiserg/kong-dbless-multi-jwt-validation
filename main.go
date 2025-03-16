package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Define the JWT generation endpoint
	router.GET("/generate-jwt/:keyfile", generateJWTHandler)

	// Start the server
	fmt.Println("JWT Generator Service started on http://localhost:8010")
	router.Run(":8010")
}

func generateJWTHandler(c *gin.Context) {
	// Get keyfile name from URL parameter
	keyFileName := c.Param("keyfile")

	// Sanitize the file path to prevent directory traversal attacks
	// Only allow .pem files from the keys directory
	if filepath.Ext(keyFileName) != ".pem" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only .pem files are allowed",
		})
		return
	}

	// Construct the key file path (from a keys directory for better security)
	keyFilePath := filepath.Join("keys", keyFileName)

	// Generate JWT
	tokenString, kid, err := generateJWT(keyFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate JWT",
			"details": err.Error(),
		})
		return
	}

	// Return the JWT token
	c.JSON(http.StatusOK, gin.H{
		"key_file": keyFileName,
		"kid":      kid,
		"jwt":      tokenString,
	})
}

// generateJWT creates a JWT signed with the specified private key file
func generateJWT(keyFilePath string) (string, string, error) {
	// Load private key from file
	privateKeyBytes, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return "", "", fmt.Errorf("error reading private key file: %v", err)
	}

	// Parse the private key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", "", fmt.Errorf("error parsing private key: %v", err)
	}

	// Generate a key ID based on SHA-256 hash of the public key
	keyID := generateKeyID(privateKey)

	// Create a token with custom headers
	token := jwt.New(jwt.SigningMethodRS256)
	
	// Set header parameters
	token.Header["kid"] = keyID    // Key ID
	token.Header["typ"] = "JWT"    // Type of token
	
	// Set standard claims
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "https://myauth.example.com"  // Issuer
	claims["sub"] = "1234567890"                  // Subject (user ID)
	claims["aud"] = "https://api.example.com"     // Audience
	claims["iat"] = time.Now().Unix()             // Issued at time
	claims["exp"] = time.Now().Add(time.Hour).Unix() // Expiration (1 hour)
	claims["nbf"] = time.Now().Unix()             // Not valid before current time
	claims["jti"] = generateUniqueID()            // JWT ID (unique identifier for this token)
	
	// Set custom claims
	claims["name"] = "John Doe"
	claims["email"] = "john.doe@example.com"
	claims["roles"] = []string{"user", "admin"}
	claims["permissions"] = map[string]bool{
		"read":  true,
		"write": true,
		"admin": true,
	}

	// Sign the token with the private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("error signing token: %v", err)
	}

	return tokenString, keyID, nil
}

// generateKeyID creates a key ID by taking the SHA-256 hash of the public key
// and encoding it as base64
func generateKeyID(privateKey *rsa.PrivateKey) string {
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalf("Error marshaling public key: %v", err)
	}
	
	hasher := sha256.New()
	hasher.Write(publicKeyDER)
	thumbprint := hasher.Sum(nil)
	
	// Return base64url encoded thumbprint as kid
	return base64.RawURLEncoding.EncodeToString(thumbprint)
}

// generateUniqueID creates a random unique ID for the JWT
func generateUniqueID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Error generating random ID: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}