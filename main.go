package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/topicai/candy"
)

const (
	defaultABACPolicyFile         = "/abac-policy.jsonl"
	defaultRestartAPIServerScript = "/restart_apiserver.sh"
	certRootPath                  = "/users"
)

// UsersBody describe http request body
// url:     http://<domain>/users
// method:  POST
type UsersBody struct {
	username  string
	namespace string
	email     string
}

func main() {
	addr := flag.String("addr", ":8080", "Listening address")
	caCrt := flag.String("ca-crt", "", "CA certificate file, in PEM format")
	caKey := flag.String("ca-key", "", "CA private key file, in PEM format")

	flag.Parse()

	if len(*caCrt) == 0 || len(*caKey) == 0 {
		glog.Fatal("Files ca.pem and ca-key.pem should be provided.")
	}
	l, e := net.Listen("tcp", *addr)
	candy.Must(e)
	// start and run the HTTP server
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", makeUsersHandler(*caCrt, *caKey))
	glog.Fatal(http.Serve(l, router))
}
func makeUsersHandler(caCrt, caKey string) http.HandlerFunc {
	return makeSafeHandler(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		candy.Must(err)

		var b UsersBody
		err = json.Unmarshal(body, &b)
		candy.Must(err)

		writeCertFiles(caCrt, caKey, certRootPath, b.username)

		// TODO: update policy
		// UpdatePolicyFile(username, namespace, defaultABACPolicyFile)
		// TODO: restart apiserver container
		// TODO: send email
		fmt.Println(body)
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
