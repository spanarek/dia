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
    "protocol": getEnv("DIAWH_PROTOCOL", "http"),
    "port": getEnv("DIAWH_PORT", "8080"),
    "cors": getEnv("DIAWH_CORS", "disabled"),
    "registry_skip_verify": getEnv("DIAWH_REGISTRY_SKIP_VERIFY", "true"),
    "ca_path": getEnv("DIAWH_CA_PATH", "./ca.pem"),
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
