# oauth2-client-test
Simple testing client for OAuth2

## Requirements
* go version 1.17+

## Usage

1. Create a JSON configuration file called "config.json" in the current directory (see below for format)
2. Run `make`
4. Run `./test` to begin the test

## Configuration File Format
Filename: config.json

Example File:

```
{
	"client_id": "client_id_from_provider",
	"client_secret": "client_secret_from_provider",
	"scopes": ["https://www.googleapis.com/auth/userinfo.email"],
	"endpoint_auth_url": "https://[oauth-provider-endpoint-authentication]",
	"endpoint_token_url": "https://[oauth-provider-endpoint-tokens]",
	"user_api_url": "https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s"
}
```

### Summary of configuration fields

* **client_id** (string)
  * Unique ID given by the provider when you setup OAuth access
* **client_secret** (string)
  * Unique token/password given by the provider when you setup OAuth access
* **scopes** (array of strings)
  * Some providers allow for custom permissions, this lets to setup the list of access scopes requested.
  * Leave empty otherwise ([])
* **endpoint_auth_url** (string)
  * URL for the endpoint of the provider which begins authentication
* **endpoint_token_url** (string)
  * URL for the endpoint of the provider which verifies token authenticity
* **user_api_url** (string)
  * URL for the API to request the info about the user with the associated auth token
  * Note: use "%s" in the URL for the placeholder where the access token will be injected
