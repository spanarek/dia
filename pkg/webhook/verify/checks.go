package verify

import (
    "archive/tar"
    "io"
    "log"
    "os"
    "crypto/x509"
    "io/ioutil"
    "encoding/pem"
    "errors"
    "strings"
    "fmt"
)


func VerifySignatureLayer(digestSliceValues string, imageDigest *string, DIAWH_CA_PATH string, uncompressedStream io.ReadCloser) error {
    tarReader := tar.NewReader(uncompressedStream)
    var tbs []byte
    // log.Print("1 Success")
    for true {
        _, err := tarReader.Next()
        if err == io.EOF {
            break
        }
        // log.Print("2 Success")
        if err != nil {
            log.Print("ExtractTar: Next() failed: %s", err.Error())
            return err
        }
        tbs, _ = ioutil.ReadAll(tarReader)
    }

    ca_cert, err := os.ReadFile(DIAWH_CA_PATH)
    if err != nil {
        log.Print("Error open CA file: ", err.Error())
        return err
    }

    err = checkCert(digestSliceValues, imageDigest, tbs, ca_cert)
    if err != nil {
        return err
    }
    return nil
}


func checkCert(digestSliceValues string, imageDigest *string, crt []byte, CA []byte) error {
  bcrt, _ := pem.Decode(crt)
  if bcrt == nil || bcrt.Type != "CERTIFICATE" {
    return errors.New("dia-image layer does not contain any CERTIFICATE in PEM format")
  }
  cert, err := x509.ParseCertificate(bcrt.Bytes)
  if err != nil {
    log.Print("failed to parse certificate: " + err.Error())
    return err
  }

  imageDigestSlice := *imageDigest
  if digestSliceValues == "" {
    digestSliceValues = "0-45"
  }
  var digestSliceStart int
	var digestSliceEnd int
	_, err = fmt.Sscanf(digestSliceValues, "%d-%d", &digestSliceStart, &digestSliceEnd)
  if err != nil {
    log.Print("failed to parse DigestSlice %d-%d: " + err.Error())
    return err
  }
  imageDigestSlice = string(imageDigestSlice[digestSliceStart:digestSliceEnd])
  if !strings.Contains(cert.Subject.CommonName, imageDigestSlice) {
    log.Print("Invalid certificate: "+cert.Subject.CommonName)
    return errors.New("Invalid certificate")
  }

  bca, _ := pem.Decode(CA)
  if bca == nil || bca.Type != "CERTIFICATE" {
    return errors.New("DIAWH_ATTESTOR_CA_CERT does not contain any certificate in PEM format")
  }
  cacert, err := x509.ParseCertificate(bca.Bytes)
  if err != nil {
    log.Print("failed to parse CA certificate: " + err.Error())
    return err
  }

  err = cert.CheckSignatureFrom(cacert)
  if err != nil {
    log.Print("Certificate signed by unknown authority "+err.Error())
    return errors.New("Certificate signed by unknown authority, "+err.Error())
  }
  return nil

}
