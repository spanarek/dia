package main

import "github.com/spanarek/dia/pkg/webhook/basic"

//Get application configuration and run server
func main() {
  //Get configuration:
  appConfig := basic.GetAppConf()
  //Running api service:
  RouterInit(appConfig)
}
