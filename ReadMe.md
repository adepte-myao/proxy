# Proxy server

This application was designed for personal education purposes.

## How to use locally:
- Clone all the files
- Open the terminal in 'proxy' folder
- Execute "go run main.go"

## How to use via Docker:
- Clone all the files
- Open the terminal in 'proxy' folder
- Execute "docker build --tag go-proxy ."
- Execute "docker run -d -p 9091:9091 --name go-proxy go-proxy"

## Simple request:
"http://localhost:9091/"

Response: "Hello from server"

## Get Links request example:
"http://localhost:9091/get-links"

Body: 
{
    "link": "https://vk.com/"
}

Response: all external links of given page