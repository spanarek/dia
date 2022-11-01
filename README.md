# DIA
## Docker image attestor

Just sign and verify your docker images with x509 certificates, keyless, really...
1. Build and push your image
2. Issue an x509 certificate: part of commonName must contain image digest(so-called digestslice)
3. Sign(push certificate) with dia-sign.sh or with your CI manually
4. Validate digestslice with gatekeeper and webhook

## Digestslice
Digestslice - a part of digest sha256. Digestslice used due to x509 commonName length limitation(64 characters).
By default digestslice = 0-45.

Example:

Full digest: 0addcc1de26ee0f660d21b01c1afdff9f59efb989331fed17334cf8a6dcd8d6b

commonName certificate must contain from 0 to 45 character: 0addcc1de26ee0f660d21b01c1afdff9f59efb98933

## Sign:

```
dia-sign.sh [ARGS] [IMAGE]
  -h, --help      print this help message
  -c, --cert      path to x509 certificate file (base64 encoded)
EXAMPLE: dia-sign.sh -c /tmp/image.crt registry.local/test-app:v1
```

## Sign with gitlab and hashicorp vault(pki):

WIP...

# Webhook
A webhook recieve your deployments and another yaml, and check digest and certificate issuer for image.
We don't check expiration dates, it's not really for images and some artifacts...

Verify request example:

```bash
curl -k -X POST --data '{ "image": "registry.local/test-app:v1",
                          "namespace": "dia-demo",
                          "imagePullSecrets": ["test-app-registry"],
                          "digestslice": "0-40" }' "http://diawh.dia.svc:8080/verify"
```
Responses:
200
```json
{"Success": "approved"}
```
500, image not signed:
```json
{"Error": "GET https://registry.local/test-app/manifests/dia-0addcc1de26ee0f660d21b01c1afdff9f59efb989331fed17334cf8a6dcd8d6b: NOT_FOUND: artifact test-app:dia-0addcc1de26ee0f660d21b01c1afdff9f59efb989331fed17334cf8a6dcd8d6b not found"}
```
403, certificate CN invalid
```json
{"Error": "Invalid certificate"}
```
403, certificate issued by unknown authority
```json
{"Error": "Certificate signed by unknown authority"}
```
403, dia image file does not contain x509 certificate on pem format
```json
{"Error": "failed to parse CA certificate"}
```

# Webhook install, helm chart:
Build an webhook image
```
docker build -t your-registry.local/diawh -f webhook.Dockerfile .
```

Webhook chart parameters: see chart/values.yaml

# Gatekeeper ConstraintTemplate

WIP...

# Arihitecture
 <img src="architecture.svg">

# TODO
Sign with JWT tokens
