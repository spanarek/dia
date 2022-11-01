package verify

import ("github.com/julienschmidt/httprouter"
        "net/http"
        "github.com/spanarek/dia/pkg/webhook/basic"
        "fmt"
        "encoding/json"
        "context"
    )

type verifyRequestBody struct {
  Image string
  Namespace string
  ImagePullSecrets []string
  DigestSlice string
}

// Post requests handler
func PostHandler() httprouter.Handle {
  return (func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    w.Header().Set("Content-Type", "application/json")
    switch location := ps.ByName("name"); location {
      default:
          jstatus := basic.JsonStatus("Error", "")
          http.Error(w, jstatus, http.StatusBadRequest)
      case "echo":
          jstatus := basic.JsonStatus("Success", "echo")
          fmt.Fprint(w, jstatus)
      case "verify":
          var vrb verifyRequestBody
          var err error
          decoder := json.NewDecoder(r.Body)
          err = decoder.Decode(&vrb)
          if err !=nil{
            jstatus := basic.JsonStatus("Error", err.Error())
            http.Error(w, jstatus, http.StatusBadRequest)
          } else {
              imageDigest, imageLayer, err := basic.GetSignImageLayer(
                context.Background(),
                vrb.Image,
                vrb.Namespace,
                vrb.ImagePullSecrets,
              )
              if err != nil {
                jstatus := basic.JsonStatus("Error", err.Error())
                http.Error(w, jstatus, http.StatusInternalServerError)
              } else {
                err := VerifySignatureLayer(
                  vrb.DigestSlice,
                  imageDigest,
                  basic.GetAppConf()["ca_path"],
                  imageLayer,
                )
                if err != nil {
                  jstatus := basic.JsonStatus("Error", err.Error())
                  http.Error(w, jstatus, http.StatusForbidden)
                } else {
                    jstatus := basic.JsonStatus("Success", "approved")
                    fmt.Fprint(w, jstatus)
                }

              }
          }
      }
  })
}
