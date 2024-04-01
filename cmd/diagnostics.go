package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func CollectDiagnostics(deploymentId string, apiKey string, essApiEndpoint string) (string, error) {

	apiURL := fmt.Sprintf("https://%s/api/v1/deployments/%s/elasticsearch/_main/diagnostics/_capture", essApiEndpoint, deploymentId)
	log.Debug("apiURL: ", apiURL)

	diagFile := fmt.Sprintf("diagnostic-%s-%s.zip", deploymentId[:6], time.Now().Format("2006-Jan-02--15_04_05"))
	log.Debug("diagFile: ", diagFile)

	// Prepare POST request
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return "", err
	}
	// Set API key in request header for authentication
	req.Header.Set("Authorization", "ApiKey "+apiKey)

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response code status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected return code: %d", resp.StatusCode)
	}

	// Read binary response
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Save response
	err = os.WriteFile(diagFile, content, 0644)
	if err != nil {
		return "", err
	}

	return diagFile, nil
}
