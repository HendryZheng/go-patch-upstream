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

		// Convert content to string
		contentStr := string(content)

		// Iterate over each mapping
		for _, mapping := range urlMapping {
			// Check if the "from" value matches the current mapping
			if val, found := mapping[from]; found {
				// Replace the target string in the content
				contentStr = strings.ReplaceAll(contentStr, val, mapping[to])
			}
		}

		// Write the updated content back to the file
		err = os.WriteFile(file, []byte(contentStr), 0644)
		if err != nil {
			fmt.Println("Error:", err)
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
