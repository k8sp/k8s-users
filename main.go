package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/k8sp/k8s-users/users"
	"github.com/topicai/candy"
)

const (
	defaultABACPolicyFile         = "./testdata/abac-policy.jsonl"
	defaultRestartAPIServerScript = "./scripts/restart_apiserver.sh"
	defaultCertFilesRootPath      = "./testdata/users"
)

func main() {
	addr := flag.String("addr", ":8080", "Listening address")
	caCrt := flag.String("ca-crt", "", "CA certificate file, in PEM format")
	caKey := flag.String("ca-key", "", "CA private key file, in PEM format")
	abacPolicyFile := flag.String("abac-policy", defaultABACPolicyFile, "Policy file with ABAC mode.")
	certFilesRootPath := flag.String("cert-root-path", defaultCertFilesRootPath, "cert files root directory")

	flag.Parse()

	if len(*caCrt) == 0 || len(*caKey) == 0 {
		glog.Fatal("Files ca.pem and ca-key.pem should be provided.")
	}

	// start and run the HTTP server
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", makeUsersHandler(*caKey, *caCrt, *certFilesRootPath, *abacPolicyFile))
	glog.Fatal(http.ListenAndServe(*addr, router))
}

func makeUsersHandler(caKey, caCrt, certFilesRootPath, abacPolicyFile string) http.HandlerFunc {
	return makeSafeHandler(func(w http.ResponseWriter, r *http.Request) {
		// load ABAC policy file
		p, err := users.LoadPoliciesfromJSONFile(abacPolicyFile)
		candy.Must(err)

		var us []users.Users

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			glog.Fatal(err)
		} else if err == io.EOF {
		}

		err = json.Unmarshal([]byte(b), &us)
		if err != nil {
			glog.Fatal(err)
		}

		for i, u := range us {

			fmt.Println(i, u.Username, u.Namespace, u.Email)

			// Update abac policy file
			if p.Exists(u) {
				p.Update(u)
			} else {
				p.Append(u)
			}

			// update cert files
			users.WriteCertFiles(caCrt, caKey, certFilesRootPath, u.Username)
		}

		p.DumpJSONFile(abacPolicyFile)

		// TODO: implement function by docker client:https://github.com/docker/docker/tree/master/client
	})
}

func makeSafeHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			}
		}()

		h(w, r)
	}
}
