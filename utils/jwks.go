package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io/fs"
	"log"
	"os"

	jose "github.com/go-jose/go-jose/v3"
)

type KeyResponse struct {
	Keys []jose.JSONWebKey `json:"keys"`
}

// generate SSH keypair
func GenerateKeypair(keySize int) ([]byte, []byte, crypto.PublicKey) {
	// generate private key file
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		log.Fatalln("Unable to generate Private Key", err)
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalln("Unable to marshal Private Key", err)
	}
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)

	// generate public key file
	publicKey := privateKey.Public()
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey.(*rsa.PublicKey))
	if err != nil {
		log.Fatalln("Unable to marshal Public Key", err)
	}
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)

	return privateKeyPem, publicKeyPem, publicKey
}

// get keyID from public key - sauce https://github.com/Azure/azure-workload-identity/blob/9893baf454a961e03ac544c42e3e30cc87998927/pkg/cmd/jwks/root.go#L180-L202
func KeyIDFromPublicKey(publicKey interface{}) string {
	publicKeyDERBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatalln("failed to serialize public key to DER format", err)
	}

	hasher := crypto.SHA256.New()
	hasher.Write(publicKeyDERBytes)
	publicKeyDERHash := hasher.Sum(nil)
	keyID := base64.RawURLEncoding.EncodeToString(publicKeyDERHash)

	return keyID
}

func GenerateJwksFromPublicKeyPem(publicKey crypto.PublicKey) []byte {
	kid := KeyIDFromPublicKey(publicKey)

	jwks := jose.JSONWebKey{
		Key:       publicKey,
		KeyID:     kid,
		Algorithm: string(jose.RS256),
		Use:       "sign",
	}

	var keys jose.JSONWebKeySet
	keys.Keys = append(keys.Keys, jwks)
	jwksJson, err := json.MarshalIndent(keys, "", "	")
	if err != nil {
		log.Fatalln("Failed to create json encoding of the JSONWebKeySet", err)
	}

	return jwksJson
}

// write to file
func WriteToFile(filepath string, content []byte, mode fs.FileMode) {
	err := os.WriteFile(filepath, content, mode)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
