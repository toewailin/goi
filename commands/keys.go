package commands

import (
	"crypto/rand"
	"crypto/rsa" // Assuming goi/utils provides PrintSuccess, etc.
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"os/exec" // For joining paths safely

	"github.com/spf13/cobra"
	// Add time package for timestamp generation
)

// Check if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// Check if OpenSSL is installed
func isOpenSSLInstalled() bool {
	cmd := exec.Command("openssl", "version")
	err := cmd.Run()
	return err == nil
}

// Generate RSA Key Pair
func generateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA private key: %v", err)
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

// Save Private Key to File
func savePrivateKeyToFile(privateKey *rsa.PrivateKey, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer file.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	err = pem.Encode(file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})
	if err != nil {
		return fmt.Errorf("failed to write private key to file: %v", err)
	}
	return nil
}

// Save Public Key to File
func savePublicKeyToFile(publicKey *rsa.PublicKey, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %v", err)
	}
	defer file.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	err = pem.Encode(file, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyBytes})
	if err != nil {
		return fmt.Errorf("failed to write public key to file: %v", err)
	}
	return nil
}

// GenerateKeys checks if files exist and generates keys if necessary
func GenerateKeys() error {
	privateKeyPath := "config/rsa_private.pem"
	publicKeyPath := "config/rsa_public.pem"

	// Check if the keys already exist
	if fileExists(privateKeyPath) && fileExists(publicKeyPath) {
		log.Println("RSA key pair already exists, skipping generation.")
		return nil
	}

	// Check if OpenSSL is installed
	if !isOpenSSLInstalled() {
		return fmt.Errorf("openssl is not installed. please install openssl to proceed")
	}

	// Generate the RSA key pair
	privateKey, publicKey, err := generateRSAKeyPair(2048)
	if err != nil {
		log.Fatalf("Error generating RSA key pair: %v", err)
	}

	// Save the private key
	err = savePrivateKeyToFile(privateKey, privateKeyPath)
	if err != nil {
		log.Fatalf("Error saving private key: %v", err)
	}

	// Save the public key
	err = savePublicKeyToFile(publicKey, publicKeyPath)
	if err != nil {
		log.Fatalf("Error saving public key: %v", err)
	}

	log.Println("RSA key pair generated and saved to files successfully!")
	return nil
}

// Define the 'generate keys' command
var GenerateKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Generate an RSA key pair",
	Long: `The 'keys' command generates a pair of RSA keys (private and public) and saves them to the config directory.

The keys will be saved as rsa_private.pem and rsa_public.pem in the config folder.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Call the GenerateKeys function to generate and save the keys
		return GenerateKeys()
	},
}

// Define the 'generate' root command
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate various resources, including RSA key pairs",
	Long: `The 'generate' command allows you to generate various resources for your project.
In this case, use the 'generate keys' command to generate RSA key pairs for your project.`,
}
