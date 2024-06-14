package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Rule struct {
	ID                string `json:"id"`
	Source            string `json:"source"`
	AzureEndpointName string `json:"azure_endpoint_name"`
	AzureEndpointGuid string `json:"azure_endpoint_guid"`
	Region            string `json:"region"`
	LinkId            string `json:"link_id"`
}

type RulesetDetails struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	IncludeByDefault bool   `json:"include_by_default"`
	Region           string `json:"region"`
	Rules            []Rule `json:"rules"`
}

type RulesetResponse struct {
	Rulesets []string `json:"rulesets"`
}

func CollectRulesets(deploymentId string, apiKey string, essApiEndpoint string) (string, error) {

	apiURL := fmt.Sprintf("https://%s/api/v1/deployments/traffic-filter/associations/deployment/%s/rulesets", essApiEndpoint, deploymentId)
	log.Debug("apiURL: ", apiURL)

	// Prepare GET request
	req, err := http.NewRequest("GET", apiURL, nil)
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

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Debug("Response Body: ", string(body))

	// Unmarshal JSON response into RulesetResponse
	var rulesetResponse RulesetResponse
	err = json.Unmarshal(body, &rulesetResponse)
	if err != nil {
		return "", err
	}

	// Extract rulesets from RulesetResponse
	rulesets := rulesetResponse.Rulesets

	// Loop over rulesets and fetch details for each ruleset
	var output interface{}
	output = make([]RulesetDetails, 0)

	for _, ruleset := range rulesets {
		// Construct URL for ruleset details
		rulesetURL := fmt.Sprintf("https://%s/api/v1/deployments/traffic-filter/rulesets/%s?include_associations=true", essApiEndpoint, ruleset)

		// Create HTTP request for ruleset details
		req, err := http.NewRequest("GET", rulesetURL, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "ApiKey "+apiKey)

		// Send HTTP request for ruleset details
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		// Read response body for ruleset details
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		log.Debug("Response Body: ", string(body))

		// Unmarshal JSON response into RulesetDetails
		var rulesetDetails RulesetDetails
		if err := json.Unmarshal(body, &rulesetDetails); err != nil {
			return "", err
		}

		output = append(output.([]RulesetDetails), rulesetDetails)
	}
	return formatRulesetOutput(output.([]RulesetDetails)), nil
}

func formatRulesetOutput(data []RulesetDetails) string {
	var ret string
	for _, ruleset := range data {
		ret += fmt.Sprintln("Ruleset ID:", ruleset.ID)
		ret += fmt.Sprintln("Name:", ruleset.Name)
		ret += fmt.Sprintln("Type:", ruleset.Type)
		ret += fmt.Sprintln("Include By Default:", ruleset.IncludeByDefault)
		ret += fmt.Sprintln("Region:", ruleset.Region)
		ret += fmt.Sprintln("Rules:")
		for _, rule := range ruleset.Rules {
			ret += fmt.Sprintln("  - ID:", rule.ID)
			if rule.Source != "" {
				ret += fmt.Sprintln("    Source:", rule.Source)
			}
			if rule.AzureEndpointName != "" {
				ret += fmt.Sprintln("    Azure Endpoint Name:", rule.AzureEndpointName)
			}
			if rule.AzureEndpointGuid != "" {
				ret += fmt.Sprintln("    Azure Endpoint Guid:", rule.AzureEndpointGuid)
			}
			if rule.Region != "" {
				ret += fmt.Sprintln("    Region:", rule.Region)
			}
			if rule.LinkId != "" {
				ret += fmt.Sprintln("    Link ID:", rule.LinkId)
			}
		}
		ret += fmt.Sprintln()
	}
	log.Debug("Ruleset output: ", ret)
	return ret
}
