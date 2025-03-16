# kong-dbless-multi-jwt-validation
This project demonstrates how to configure Kong API Gateway to validate multiple JWT tokens, where each token is signed with a different private key. The setup ensures secure authentication by dynamically verifying JWT signatures using their respective public keys.

## Run the project

To start the project, run:


```sh
docker-compose -f docker-compose-jwt.yml up --build
```
This will build and launch the necessary services.


## Getting the JWT tokens

To generate different JWT tokens, use the following endpoints:



```sh
http://localhost:8000/generate-jwt/private_key_1.pem
http://localhost:8000/generate-jwt/private_key_2.pem
http://localhost:8000/generate-jwt/private_key_3.pem
```
Each request will return a JWT signed with the corresponding private key stored in the keys folder.


The response format is:

```sh
{
jwt: "eyJhbGciOiJSUzI1NiIsImtpZCI6IkNBeTU2QlpseGNtbkxQazJoRVhiZHZXalpMOW9pMkJ3U0xHNmo5TmxSMEUiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJodHRwczovL2FwaS5leGFtcGxlLmNvbSIsImVtYWlsIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJleHAiOjE3NDIwOTI4NDEsImlhdCI6MTc0MjA4OTI0MSwiaXNzIjoiaHR0cHM6Ly9teWF1dGguZXhhbXBsZS5jb20iLCJqdGkiOiJfYUtBOGszVFllLTZKdmpUS0lxZDVRIiwibmFtZSI6IkpvaG4gRG9lIiwibmJmIjoxNzQyMDg5MjQxLCJwZXJtaXNzaW9ucyI6eyJhZG1pbiI6dHJ1ZSwicmVhZCI6dHJ1ZSwid3JpdGUiOnRydWV9LCJyb2xlcyI6WyJ1c2VyIiwiYWRtaW4iXSwic3ViIjoiMTIzNDU2Nzg5MCJ9.fAPgjDBgvbVmU_SIiXj6u3P1MkmX_7hRzM3uPfg6sb-EEvxvE3GBea0EduM4JtAgFqB7LZCHAHmtKdvMJhEVMmBzwmhilNqXwbWCWIIQmDNqFMPjpoRnaYDHvv6qE6Z2jOGT2tKCGDbCN6pxQsTVp5dUEUjEJfDWuTMUHNEDehQVBaQgpmsFqxmLE3-g5eZ01XlzdIbs_WxL9jN2Eq5_iThvEPh1aiVHb7a-KHm54-D8_rd-oLVhk5zpBNXDOWjWdyyzw4oCh3lvW4Y7ErSxHfWbUpKlyDPnABcJ5AWEtmP_mNT8Kbz7f5V2F3Hdx6OoXIa3vAbtnHYPyGMhPc2i8Q",
key_file: "private_key_3.pem",
kid: "CAy56BZlxcmnLPk2hEXbdvWjZL9oi2BwSLG6j9NlR0E"
}
```

## Validate the JWT

To examine the structure and content of your JWT, visit https://jwt.io and paste your token into the decoder. This tool will display the header, payload, and verify the signature of your JWT.


To validate a JWT, send a request to Kong with the generated token:


```sh
curl -X GET "http://localhost:8000/echo/get" \
     -H "Authorization: Bearer YOUR_GENERATED_JWT" -v
```
__Note__: Replace YOUR_GENERATED_JWT with the actual token received from the previous step.

## How It Works

Kong is configured to validate JWTs using different public keys.

JWTs are signed with private keys stored in the keys folder.

The Kong JWT plugin verifies the token against the correct public key based on the Key ID (kid) in the JWT header.

If the token is valid, the request is forwarded to the upstream service. Otherwise, Kong rejects it with an authentication error.


## Command to generate keys


```sh
openssl genrsa -out private_key.pem 2048; openssl rsa -in private_key.pem -pubout -out public_key.pem
```