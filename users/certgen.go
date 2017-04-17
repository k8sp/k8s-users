package users

import (
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/k8sp/sextant/golang/certgen"
	"github.com/topicai/candy"
	"github.com/zh794390558/go-study/cert/gencrt"
)

// genCerts  generate key and crt files
func genCerts(caCrt, caKey, username string) ([]byte, []byte) {
	out, e := ioutil.TempDir("", "")
	candy.Must(e)

	defer func() {
		if e = os.RemoveAll(out); e != nil {
			log.Printf("certgen.Gen failed deleting %s", out)
		}
	}()

	key := path.Join(out, username+"-key.pem")
	csr := path.Join(out, username+"-csr.pem")
	crt := path.Join(out, username+"-crt.pem")

	//openssl genrsa -out <username>-key.pem 2048
	//openssl req -new -key <username>-key.pem -out <username>.csr -subj "/CN=$1"
	//openssl x509 -req -in <username>.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out users/$1/$1.pem -days 365
	certgen.Run("openssl", "genrsa", "-out", key, "2048")
	certgen.Run("openssl", "req", "-new", "-key", key, "-out", csr, "-subj",
		"/CN="+username)
	certgen.Run("openssl", "x509", "-req", "-in", csr, "-CA", caCrt, "-CAkey",
		caKey, "-CAcreateserial", "-out", crt, "-days", "365")

	k, e := ioutil.ReadFile(key)
	candy.Must(e)

	c, e := ioutil.ReadFile(crt)
	candy.Must(e)

	return k, c
}

func genUserCert(caCrt, caKey, userName string) (key []byte, crt []byte) {
	var cacrt *x509.Certificate
	var cakey *rsa.PrivateKey
	var err error

	//load caCrt
	cacrt, err = gencrt.ParseCertificate(caCrt)
	candy.Must(err)

	//load caKey
	cakey, err = gencrt.ParseRSAPrivateKey(caKey)
	candy.Must(err)

	//gen key
	var priv interface{}
	priv, err = gencrt.GenerateRSAPrivKey(2048)
	candy.Must(err)
	key = gencrt.PemEncodeToMemory(gencrt.PemBlockForKey(priv))

	//gen crt
	derBytes := gencrt.CreateUserCertificate(cacrt, cakey, priv, "testUser", 24*time.Hour)
	crt = gencrt.PemEncodeToMemory(gencrt.PemBlockForCrt(derBytes))

	return
}

//WriteCertFiles generate cert files in #certRootPath
func WriteCertFiles(caCrt, caKey, certRootPath, username string) (crtFile, keyFile string) {
	userPath := path.Join(certRootPath, username)

	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		os.MkdirAll(userPath, 0744)
	}

	//key, crt := genCerts(caCrt, caKey, username)
	key, crt := genUserCert(caCrt, caKey, username)

	crtFile = path.Join(userPath, username+"-crt.pem")
	keyFile = path.Join(userPath, username+"-key.pem")

	err := ioutil.WriteFile(crtFile, crt, 0644)
	candy.Must(err)

	err = ioutil.WriteFile(keyFile, key, 0644)
	candy.Must(err)

	return
}
