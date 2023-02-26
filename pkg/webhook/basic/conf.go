//Working whith application configurations
package basic

import ("os"
        // "log"
      )

// getEnv get key environment variable if exist otherwise return defalutValue
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return defaultValue
    }
    return value
}

func GetAppConf() map[string]string {
  configs := map[string]string{
    "protocol": getEnv("DIAWH_PROTOCOL", "https"),
    "port": getEnv("DIAWH_PORT", "8080"),
    "cors": getEnv("DIAWH_CORS", "disabled"),
    "registry_skip_verify": getEnv("DIAWH_REGISTRY_SKIP_VERIFY", "true"),
    "attestor_ca_path": getEnv("DIAWH_ATTESTOR_CA_CERT", "/etc/pki/tls/certs/ca.pem"),
    "tls_cert": getEnv("DIAWH_TLS_CERT", "/etc/pki/tls/certs/diawh.crt"),
    "tls_key": getEnv("DIAWH_TLS_KEY", "/etc/pki/tls/private/diawh.key"),
    "digest_slice": getEnv("DIGEST_SLICE", "0-45"),
  }
  return configs
}

//Create json response status messages from simple string
func JsonStatus(status, msg string) string{
  switch status {
  case "Error":
    return "{\"Error\": \""+msg+"\"}"
  case "Success":
    return "{\"Success\": \""+msg+"\"}"
  default:
    return "{\"Info\": \""+msg+"\"}"
  }
}
