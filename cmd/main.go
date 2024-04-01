package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultApiEndpoint = "api.elastic-cloud.com"
	DeploymentIdRegexp = `^[a-f0-9]{32}$`
	OptA               = "main-elasticsearch"
	OptB               = "elasticsearch"
	NotFound           = "resource_not_found"
)

func main() {
	var verboseMode, commandMode, diagMode, tfMode bool
	var response string
	var err error

	flag.BoolVar(&verboseMode, "v", false, "Verbose")
	flag.BoolVar(&diagMode, "d", false, "Diagnostics mode")
	flag.BoolVar(&tfMode, "t", false, "Traffic filters mode")
	flag.Parse()

	if verboseMode {
		log.SetLevel(log.DebugLevel)
	}

	// Get API key from environment variable
	apiKey := os.Getenv("ELASTIC_ESS_KEY")
	if apiKey == "" {
		log.Fatal("ELASTIC_ESS_KEY environment variable not set")
	}
	log.Debug("ELASTIC_ESS_KEY: ", apiKey)

	// Get Elasticsearch Service API endpoint from environment variable
	// or assign the default one
	essApiEndpoint := os.Getenv("ELASTIC_ESS_APIENDPOINT")
	if essApiEndpoint == "" {
		essApiEndpoint = DefaultApiEndpoint
	}
	log.Debug("essApiEndpoint: ", essApiEndpoint)

	// Get ESS deployment id from program argument
	if len(flag.Args()) < 1 {
		log.Fatal("Usage: ess [-v] [-d|-t] <deployment_id> [Elasticsearch GET command]")
	}
	deploymentId := flag.Arg(0)
	command := flag.Arg(1)

	// If not running in Traffic Filters mode or Diagnostics Mode, assume we are running in commandMode
	if !tfMode && !diagMode {
		commandMode = true
		if command == "" {
			command = "/"
		}
	}

	// If command is present, we are running in commandMode
	// Ensure command starts with a '/'
	if command != "" {
		commandMode = true
		if command[0] != '/' {
			command = "/" + command
		}
	}

	// Can't run certain modes together
	if commandMode && (tfMode || diagMode) {
		log.Fatal("Command and Diagnostics / Traffic Filters modes are mutually exclusive.")
	}

	// Validate deployment id format
	pattern := regexp.MustCompile(DeploymentIdRegexp)
	if !pattern.MatchString(deploymentId) {
		log.Fatal("Invalid deployment_id format: ", deploymentId)
	}

	// Diagnostics mode
	if diagMode {
		log.Info("Collecting diagnostics...")
		filename, err := CollectDiagnostics(deploymentId, apiKey, essApiEndpoint)
		if err != nil {
			log.Error(err.Error())
		} else {
			log.Info("Done")
			log.Info("Diagnostics saved in ", filename)
		}
	}

	// Traffic Filters Mode
	if tfMode {
		log.Info("Collecting ruleset...")
		output, err := CollectRulesets(deploymentId, apiKey, essApiEndpoint)
		if err != nil {
			log.Error(err.Error())
		} else {
			log.Info("Done")
			fmt.Print(output)
		}
	}

	// Command Mode
	if commandMode {
		response, err = runCommand(deploymentId, apiKey, essApiEndpoint, OptA, command)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Resend command with a different URL
		if strings.Contains(response, NotFound) {
			response, err = runCommand(deploymentId, apiKey, essApiEndpoint, OptB, command)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		fmt.Println(response)
	}
}

func runCommand(deploymentId string, apiKey string, essApiEndpoint string, opt string, command string) (string, error) {
	apiURL := fmt.Sprintf("https://%s/api/v1/deployments/%s/elasticsearch/%s/proxy%s", essApiEndpoint, deploymentId, opt, command)
	log.Debug("apiURL: ", apiURL)

	// Prepare GET request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	// Set API key in request header for authentication
	req.Header.Set("Authorization", "ApiKey "+apiKey)
	// Set Managment Header
	req.Header.Set("X-Management-Request", "true ")

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Debug("Response Body: ", string(body))

	return string(body), nil
}
