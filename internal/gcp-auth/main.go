package gcpauth

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

const (
	GCP_PROJECT_NUMBER        = "394216874503"
	GCP_PROJECT_ID            = "trivago-415510"
	WIF_POOL_ID               = "test-kind"
	WIF_PROVIDER_ID           = "test-kind"
	K8S_SA_TOKEN_PATH         = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	TOKEN_SCOPE               = "https://www.googleapis.com/auth/cloud-platform"
	GCP_SERVICE_ACCOUNT_EMAIL = "use-me-to-impersonate-and-stor@trivago-415510.iam.gserviceaccount.com"
)

func MakeSTSRequest() (string, error) {
	url := "https://sts.googleapis.com/v1/token"

	tokenBytes, err := os.ReadFile(K8S_SA_TOKEN_PATH)
	if err != nil {
		fmt.Errorf("failed to read service account token: %w", err)
	}

	token := string(tokenBytes)

	// Create empty request body for dummy request
	dummyJSON := fmt.Sprintf(`{
		"grant_type": "urn:ietf:params:oauth:grant-type:token-exchange",
		"audience": "//iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/providers/%s",
		"scope": "%s",
		"requested_token_type": "urn:ietf:params:oauth:token-type:access_token",
		"subject_token": "%s",
		"subject_token_type": "urn:ietf:params:oauth:token-type:jwt"
	}`, GCP_PROJECT_NUMBER, WIF_POOL_ID, WIF_PROVIDER_ID, TOKEN_SCOPE, token)

	body := bytes.NewBuffer([]byte(dummyJSON))

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("STS response status: %s\n", resp.Status)
	respBody := &bytes.Buffer{}
	_, err = respBody.ReadFrom(resp.Body)

	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("STS Federated Token : %s\n", respBody.String())
	return resp.Status, nil

	// TODO : CALL IAM API to impersonate the service account

}
