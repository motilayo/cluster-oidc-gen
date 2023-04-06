package utils

import (
	"os"
	"encoding/json"
	"strings"
)

type Discovery struct {
	Issuer                  string   `json:"issuer"`
	JwksUri                 string   `json:"jwks_uri"`
	AuthorizationEndpoint   string   `json:"authorization_endpoint"`
	ResponseTypesSupported  []string `json:"response_types_supported"`
	SubjectTypesSupported   []string `json:"subject_types_supported"`
	IdTokenSigningAlgValues []string `json:"id_token_signing_alg_values_supported"`
	ClaimsSupported         []string `json:"claims_supported"`
}

func CreateDiscoveryJson(oidcDiscoveryFileTemplate string, oidcProviderURL string) []byte {
	// Read the template JSON file
	tplData, err := os.ReadFile(oidcDiscoveryFileTemplate)
	if err != nil {
		panic(err)
	}

	// Convert the template JSON data to a string
	tplStr := string(tplData)

	// Replace the placeholder values with actual values
	tplStr = strings.ReplaceAll(tplStr, "$OIDC_PROVIDER", oidcProviderURL)

	// Unmarshal the JSON data into a Discovery struct
	var discovery Discovery
	err = json.Unmarshal([]byte(tplStr), &discovery)
	if err != nil {
		panic(err)
	}

	// Update some of the fields in the Discovery struct
	discovery.JwksUri = "https://" + oidcProviderURL + "/keys.json"
	discovery.ResponseTypesSupported = []string{"id_token", "code id_token"}
	discovery.IdTokenSigningAlgValues = []string{"RS256", "HS256"}

	// Marshal the updated Discovery struct back into JSON data
	newData, err := json.MarshalIndent(discovery, "", "    ")
	if err != nil {
		panic(err)
	}

	return newData
}
