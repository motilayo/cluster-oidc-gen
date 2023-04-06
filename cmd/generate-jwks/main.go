package main

import (
	"flag"

	"github.com/motilayo/cluster-oidc-gen/utils"
	log "github.com/sirupsen/logrus"
)

const (
	privateKeyFilePath        = "oidc/sa-signer.key"
	pkcs8PublicKeyFilePath    = "oidc/sa-signer-pkcs8.pub"
	jwksFilePath              = "oidc/keys.json"
	oidcDiscoveryFilePath     = "oidc/discovery.json"
	oidcDiscoveryFileTemplate = "discovery.json.tpl"
	fileMode                  = 0600
)

func main() {
	bucketName := flag.String("bucket-name", "", "Name of AWS S3 bucket")
	flag.Parse()

	log.Info("Generating SSH key pair")
	privateKeyPem, publicKeyPem, publicKey := utils.GenerateKeypair(4096)
	utils.WriteToFile(privateKeyFilePath, privateKeyPem, fileMode)
	utils.WriteToFile(pkcs8PublicKeyFilePath, publicKeyPem, fileMode)
	jwks := utils.GenerateJwksFromPublicKeyPem(publicKey)
	utils.WriteToFile(jwksFilePath, jwks, fileMode)

	//create s3 bucket
	bucketBasics := utils.CreateBucketBasics()

	bucketBasics.CreateBucket(*bucketName)
	oidcProviderURL := utils.GenerateS3URL(*bucketName)
	discoveryJson := utils.CreateDiscoveryJson(oidcDiscoveryFileTemplate, oidcProviderURL)
	utils.WriteToFile(oidcDiscoveryFilePath, discoveryJson, fileMode)
	bucketBasics.UploadToS3(*bucketName, jwksFilePath, "keys.json")
	bucketBasics.UploadToS3(*bucketName, oidcDiscoveryFilePath, ".well-known/openid-configuration")

}
