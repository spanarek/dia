package verify

import ("github.com/julienschmidt/httprouter"
        "net/http"
        "github.com/spanarek/dia/pkg/webhook/basic"
        "fmt"
        "encoding/json"
        "context"
        "k8s.io/api/admission/v1beta1"
        corev1 "k8s.io/api/core/v1"
        "log"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "regexp"
    )

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
        arReview := v1beta1.AdmissionReview{}
        err := json.NewDecoder(r.Body).Decode(&arReview)
          if err != nil {
            log.Print("Cannot decode AdmissionReview json: ", err.Error())
            jstatus := basic.JsonStatus("Error", err.Error())
        		http.Error(w, jstatus, http.StatusBadRequest)
        		return
        	} else if arReview.Request == nil {
            log.Print("Invalid request body")
            jstatus := basic.JsonStatus("Error", "invalid request body")
        		http.Error(w, jstatus, http.StatusBadRequest)
        		return
        	}
         arReview.Response = &v1beta1.AdmissionResponse{
            UID:     arReview.Request.UID,
            Allowed: true,
          }
          denyReview := func (err error) {
            arReview.Response.Allowed = false
            arReview.Response.Result = &metav1.Status{
              Message: err.Error(),
            }
            json.NewEncoder(w).Encode(&arReview)
          }
          ns := arReview.Request.Namespace
          pod := corev1.Pod{}
          err = json.Unmarshal(arReview.Request.Object.Raw, &pod)
          if err != nil {
             log.Print("Cannot decode review object, not a pod? namespace: ", ns,
                         ", object: ", arReview.Request.Name)
             denyReview(err)
             return
          }
          var imagePullSecrets []string
          for _, imagePullSecret := range pod.Spec.ImagePullSecrets {
            imagePullSecrets = append(imagePullSecrets, imagePullSecret.Name)
          }
          log.Print("Start review containers for namespace: ", ns,
                      ", object: ", arReview.Request.Name)
          allContainers := append(pod.Spec.Containers, pod.Spec.InitContainers...)
          ctx := context.Background()
          for _, container := range allContainers {
                imageRepoWithDigest, _ := regexp.MatchString("@[a-z0-9]+([+._-][a-z0-9]+)*:[a-zA-Z0-9=_-]+", container.Image)
                if imageRepoWithDigest {
                  log.Print("Skip container, deployed with digest: ", container.Image)
                  continue
                }
                imageDigest, imageLayer, err := basic.GetSignImageLayer(
                  ctx,
                  container.Image,
                  ns,
                  imagePullSecrets,
                )
                if err != nil {
                  log.Print("Error getting Sign layer: ", err.Error())
                  denyReview(err)
                  return
                } else {
                  err := VerifySignatureLayer(
                    basic.GetAppConf()["digest_slice"],
                    imageDigest,
                    basic.GetAppConf()["attestor_ca_path"],
                    imageLayer,
                  )
                  if err != nil {
                    log.Print("Error verify signature: ", err.Error())
                    denyReview(err)
                    return
                  }
                }
              }
            json.NewEncoder(w).Encode(&arReview)
      }
  })
}
