### Simple API designed for checking, if resource is allowing requests from your location

This is simple Golang web server with only one endpoint. When it receives request, it sends HTTP request to target resource and processes response
 + If response code is 403, we are sure that resource is blocking access. 
 + If response contains phrases like "Access denied under U.S. Export Administration Regulations", we can assume that resource is blocking access.

The main idea is to deploy instances of this server in different countries, then make requests to all of them to decide if resource is really blocking access only for users from particular countries.

*NOTE: This is very simple application. Web pages, are not rendered, so if error message is displayed in iframe or by JS code, it can be ignored. However, 403 status code saves the day sometimes.*

Configuration parameters are passed using environment variables:

+ HOST and PORT - host and port for web server (default to `localhost` and `8080`)
+ MAX_CONTENT_LENGTH - maximum Content-Length of response of the resource. If response is bigger, error 500 is returned (defaults to 10 000 000 bytes)
+ PHRASES_FILE - path to custom file with phrases, that will indicate that page is 403 error page. File must be just plain text, phrases separated by \n (defaults to `default_phrases.txt`)
+ FILTER_TAGS - comma-separated (e.g. `img,script,meta`) list of HTML tags to ignore (all text inside such tags will not be searched for phrases). Defaults to `script`

[See OpenAPI spec](https://captaindno.github.io/access-tester/)