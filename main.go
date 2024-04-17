package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Check if there are enough command-line arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go from=staging to=production")
		os.Exit(1)
	}

	args := os.Args[1:]

	// Parse the parameters
	params := make(map[string]string)
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			fmt.Println("Invalid parameter format:", arg)
			os.Exit(1)
		}
		params[parts[0]] = parts[1]
	}

	// Get the "from" and "to" parameters
	from, ok := params["from"]
	if !ok {
		fmt.Println("Missing 'from' parameter")
		os.Exit(1)
	}
	to, ok := params["to"]
	if !ok {
		fmt.Println("Missing 'to' parameter")
		os.Exit(1)
	}

	// Load URL mapping from JSON file
	urlMapping, err := loadURLMapping("url_mapping.json")
	if err != nil {
		fmt.Println("Error loading URL mapping:", err)
		os.Exit(1)
	}

	// Get the list of JSON files in the current directory
	files, err := filepath.Glob("tmp/api*.json")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Iterate over each JSON file
	for _, file := range files {
		fmt.Println("Processing file:", file)

		// Read the content of the JSON file
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Unmarshal JSON content into a map
		var data map[string]interface{}
		if err := json.Unmarshal(content, &data); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue
		}

		// Get the target URL from the JSON data
		targetURL, ok := getTargetURL(data)
		if !ok {
			fmt.Println("Key 'api_definition.proxy.target_url' not found in", file)
			continue
		}

		// Check if the targetURL is a string
		url, ok := targetURL.(string)
		if !ok {
			fmt.Println("Target URL is not a string in", file)
			continue
		}

		// Replace the URL with the value from the URL mapping
		for _, mapping := range urlMapping {
			if val, found := mapping[from]; found && val == url {
				url = val
				if newVal, found := mapping[to]; found {
					url = newVal
				}
				break
			}
		}

		// Update the target URL in the JSON data
		data["api_definition"].(map[string]interface{})["proxy"].(map[string]interface{})["target_url"] = url

		// Marshal the updated data back to JSON format
		updatedContent, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			continue
		}

		// Write the updated content back to the file
		if err := os.WriteFile(file, updatedContent, 0644); err != nil {
			fmt.Println("Error writing to file:", err)
			continue
		}

		fmt.Println("Successfully updated", file)
	}
}

// loadURLMapping loads URL mapping from the specified JSON file
func loadURLMapping(filename string) ([]map[string]string, error) {
	// Read the content of the JSON file
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON content into a slice of maps
	var mappings []map[string]string
	if err := json.Unmarshal(content, &mappings); err != nil {
		return nil, err
	}

	return mappings, nil
}

// getTargetURL finds and returns the target URL from the JSON data
func getTargetURL(data map[string]interface{}) (interface{}, bool) {
	proxy, ok := data["api_definition"].(map[string]interface{})["proxy"].(map[string]interface{})
	if !ok {
		return nil, false
	}

	targetURL, ok := proxy["target_url"]
	return targetURL, ok
}
