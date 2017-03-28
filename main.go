package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/k8sp/k8s-users/users"
	"github.com/topicai/candy"
)

const (
	defaultABACPolicyFile         = "./abac_policy.jsonl"
	defaultRestartAPIServerScript = "./restart_apiserver.sh"
	defaultCertFilesRootPath      = "./users"
)

// UsersBody describe http request body
// url:     http://<domain>/users
// method:  POST
type UsersBody struct {
	Username  string
	Namespace string
	Email     string
}

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
	l, e := net.Listen("tcp", *addr)
	candy.Must(e)
	// start and run the HTTP server
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", makeUsersHandler(*caKey, *caCrt, *certFilesRootPath, *abacPolicyFile))
	glog.Fatal(http.Serve(l, router))
}
func makeUsersHandler(caKey, caCrt, certFilesRootPath, abacPolicyFile string) http.HandlerFunc {
	return makeSafeHandler(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var u UsersBody
		err := decoder.Decode(&u)
		candy.Must(err)

		users.WriteCertFiles(caKey, caCrt, certFilesRootPath, u.Username)

		users.UpdatePolicyFile(u.Username, u.Namespace, abacPolicyFile)
		// TODO: implement function by docker client:https://github.com/docker/docker/tree/master/client
		users.RestartContainerByKeyword("apiserver")
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
