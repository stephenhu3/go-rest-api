package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "log"
    "math/big"
    "net"
    "time"
    "errors"
    "os"
)

// helper function to create a cert template with a serial number and other required fields
func CertTemplate() (*x509.Certificate, error) {
    // generate a random serial number (a real cert authority would have some logic behind this)
    serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
    serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
    if err != nil {
        return nil, errors.New("failed to generate serial number: " + err.Error())
    }

    tmpl := x509.Certificate{
        SerialNumber:          serialNumber,
        Subject:               pkix.Name{Organization: []string{"Yhat, Inc."}},
        SignatureAlgorithm:    x509.SHA256WithRSA,
        NotBefore:             time.Now(),
        NotAfter:              time.Now().Add(time.Hour), // valid for an hour
        BasicConstraintsValid: true,
    }
    return &tmpl, nil
}

func CreateCert(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) ( err error) {
	
	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
   
    certOut, err := os.Create("cert.pem")
    if err != nil {
     log.Fatalf("failed to open cert.pem for writing: %s", err)
    }
    pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
    certOut.Close()
    return
}

func generateCertKeyPEM(){
    // generate a new private key-pair
    rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        log.Fatalf("generating random key: %v", err)
    }

    rootCertTmpl, err := CertTemplate()
    if err != nil {
        log.Fatalf("creating cert template: %v", err)
    }

    // describe what the certificate will be used for
    rootCertTmpl.IsCA = true
    rootCertTmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
    rootCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
    rootCertTmpl.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}

    //Create roote certificate and write to cert.pem file
    if CreateCert(rootCertTmpl, rootCertTmpl, &rootKey.PublicKey, rootKey) != nil {
        log.Fatalf("error creating cert: %v", err)
    }

    keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        log.Print("failed to open key.pem for writing:", err)
        return
    }
    // pem.Encode(keyOut, pemBlockForKey(rootKey))
    pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rootKey), })
    keyOut.Close()
}