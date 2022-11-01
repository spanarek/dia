//Http request handlers package
package main

import ("github.com/julienschmidt/httprouter"
        "net/http"
        "log"
        "github.com/spanarek/dia/pkg/webhook/verify"
    )

//Initialization http server
func RouterInit(appConfig map[string]string) {
  //Initialization requests config:
  router := httprouter.New()
  basicPath := "/"
  if appConfig["cors"] == "enabled" {
      router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers
        header := w.Header()
        header.Set("Access-Control-Allow-Methods", "POST, OPTIONS, HEAD")
        header.Set("Access-Control-Allow-Origin", "*")
        header.Set("Access-Control-Allow-Headers", "Content-Type")
      // Adjust status code to 204
      w.WriteHeader(http.StatusNoContent)
     })
  }

  router.POST(basicPath+":name", verify.PostHandler())

  log.Print("DIA webhook started on port "+appConfig["port"])

  //Running service:
  switch appConfig["protocol"] {
  default:
    log.Fatal("Unsupported api protocol")
  case "http":
    log.Fatal(http.ListenAndServe(
    ":"+appConfig["port"],
    router))
  case "https":
    log.Fatal(http.ListenAndServeTLS(
      ":"+appConfig["port"],
      "/etc/ssl/certs/diawh.crt",
      "/etc/ssl/private/diawh.key",
      router))
  }
}
