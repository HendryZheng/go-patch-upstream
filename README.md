# Tyk Sync
For more details on tyk sync navigate to here
https://tyk.io/docs/tyk-sync/#example-transfer-from-one-tyk-dashboard-to-another

After tyk-sync completed, you should be able to see list of apis like below

<img width="365" alt="image" src="https://github.com/HendryZheng/go-patch-upstream/assets/2715449/1ecf0538-62a0-46aa-acb4-81bedf2965a6">

The script will look through tmp/api*.json and look for `api_definition.proxy.target_url` and replace with each env accordingly

# URL Replacer After Tyk Sync

This Go script is designed to replace URLs in JSON files based on a mapping provided in a separate JSON file.

## Usage

1. Ensure you have Go installed on your system.
2. Clone or download this repository.
3. Navigate to the directory containing the Go script and the `url_mapping.json` file.
4. Run the script using the `go run` command, providing the necessary parameters:

```
go run main.go from=staging to=production
```


Replace `staging` and `production` with the appropriate values for your use case.

## Parameters

- **from**: Specifies the source environment or URL to be replaced.
- **to**: Specifies the target environment or URL to replace the source with.

## JSON Files

- **url_mapping.json**: Contains the mapping of URLs from one environment to another. Each object in the array represents a set of mappings.

Example:
```json
[
   {
       "dev": "http://upstream-dev.com",
       "staging": "http://upstream-staging.com",
       "production": "http://upstream-production.com"
   },
   {
       "dev": "http://aa-upstream-dev.com",
       "staging": "http://aa-upstream-staging.com",
       "production": "http://aa-upstream-production.com"
   }
]
