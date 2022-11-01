package basic

import (
  kubernetesConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
  "k8s.io/client-go/kubernetes"
  "log"
  "github.com/google/go-containerregistry/pkg/authn/k8schain"
  "github.com/google/go-containerregistry/pkg/v1/remote"
  "github.com/google/go-containerregistry/pkg/name"
  "io"
  "context"
  "net/http"
  "crypto/tls"
)


func newK8SClient() (kubernetes.Interface, error) {
	kubeConfig, err := kubernetesConfig.GetConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(kubeConfig)
}

type diaImage struct {
  origDigest string
  signImageRef string
}

func GetSignImageLayer(ctx context.Context, origImage, ns string, imagePullSecrets []string) (*string, io.ReadCloser, error) {
  k8sClient, err := newK8SClient()
  if err != nil {
    log.Print("error creating k8s client: %s", err)
  }

  chainOpts := k8schain.Options{
    Namespace:          ns,
    ImagePullSecrets: imagePullSecrets,
  }


  authChain, err := k8schain.New(
    ctx,
    k8sClient,
    chainOpts,
  )
  if err != nil {
    log.Print(err, ": failed to create k8schain authentication, opts: ", chainOpts)
    return nil, nil, err
  }

  options := []remote.Option{
    remote.WithAuthFromKeychain(authChain),
  }

  if GetAppConf()["registry_skip_verify"] == "true" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // nolint:gosec
		}
		options = append(options, remote.WithTransport(tr))
    // log.Print("insecure")
	}


 // examples
 // origImage := "registry.local/test/demo-attestor:v01"
 // diaImage := "registry.local/test/demo-attestor:dia-c71fc07fbc17a5a49075f2db72e32b3b460331425b0049916cdc00071e0933c6"
 diaImage, err := getDiaImageInfo(origImage, options)
 if err != nil {
   return nil, nil, err
 }

 ref, err := name.ParseReference(diaImage.signImageRef)
 if err != nil {
   return nil, nil, err
 }

 descriptor, err := remote.Get(ref, options...)
 if err != nil {
   return nil, nil, err
 }

 image, err := descriptor.Image()
 if err != nil {
   return nil, nil, err
 }

 configFile, err := image.Manifest()
 if err != nil {
   return nil, nil, err
 }

 diaLayer, err := image.LayerByDigest(configFile.Layers[0].Digest)
 // compressedIm, err := configFile.Layers[0].Compressed
 if err != nil {
   return nil, nil, err
 }


 uncompressedLayerContent, err := diaLayer.Uncompressed()
 if err != nil {
   return nil, nil, err
 }

 return &diaImage.origDigest, uncompressedLayerContent, nil

}

func getDiaImageInfo(origImageRef string, options []remote.Option) (*diaImage, error) {
  ref, err := name.ParseReference(origImageRef)
  if err != nil {
    return nil, err
  }

  descriptor, err := remote.Get(ref, options...)
  if err != nil {
    return nil, err
  }

  image, err := descriptor.Image()
  if err != nil {
    return nil, err
  }

  imageDigest, err := image.Digest()
  if err != nil {
    return nil, err
  }
  repo := ref.Context()
  // log.Print()
  signImageRef := repo.Name()+":dia-"+imageDigest.Hex
  di := diaImage{
    origDigest: imageDigest.Hex,
    signImageRef: signImageRef,
  }

  return &di, nil

}
