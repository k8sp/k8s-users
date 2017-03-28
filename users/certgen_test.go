package users

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/k8sp/sextant/golang/certgen"
	"github.com/stretchr/testify/assert"
	"github.com/topicai/candy"
)

func TestGenCert(t *testing.T) {
	out, e := ioutil.TempDir("", "")
	candy.Must(e)
	defer func() {
		if e = os.RemoveAll(out); e != nil {
			log.Printf("Generator.Gen failed deleting %s", out)
		}
	}()
	caKey, caCrt := certgen.GenerateRootCA(out)
	key, crt := genCerts(caKey, caCrt, "test")
	assert.True(t, strings.HasPrefix(string(key), "-----BEGIN RSA PRIVATE KEY-----"))
	assert.True(t, strings.HasSuffix(string(key), "-----END RSA PRIVATE KEY-----\n"))

	assert.True(t, strings.HasPrefix(string(crt), "-----BEGIN CERTIFICATE-----"))
	assert.True(t, strings.HasSuffix(string(crt), "-----END CERTIFICATE-----\n"))
}
