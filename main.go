package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"

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
	adminEmail := flag.String("admin-email", "", "admin's email to send crt and key for users, like: admin@domain.com")
	adminSecrt := flag.String("admin-secrt", "", "admin's secrt to send crt and key for users")
	smtpsrv := flag.String("smtp-svc-addr", "", "SMTP server address with port")

	flag.Parse()

	if len(*caCrt) == 0 || len(*caKey) == 0 || len(*adminEmail) == 0 || len(*adminSecrt) == 0 || len(*smtpsrv) == 0 {
		glog.Fatal("Files ca.pem , ca-key.pem , admin email and secrt, smtp server address should be provided.")
	}

	smtp := users.NewSmtpInfo(*smtpsrv, *adminEmail, *adminSecrt)

	// start and run the HTTP server
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", makeUsersHandler(*caKey, *caCrt, *certFilesRootPath, *abacPolicyFile, smtp))
	glog.Fatal(http.ListenAndServe(*addr, router))
}

func makeUsersHandler(caKey, caCrt, certFilesRootPath, abacPolicyFile string, smtp *users.SmtpInfo) http.HandlerFunc {
	return makeSafeHandler(func(w http.ResponseWriter, r *http.Request) {
		// smtp message pool
		go smtp.SMTPSvcPool()

		// load ABAC policy file
		p, err := users.LoadPoliciesfromJSONFile(abacPolicyFile)
		candy.Must(err)

		//gen Policy for user
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
			crtFile, keyFile := users.WriteCertFiles(caCrt, caKey, certFilesRootPath, u.Username)

			// send email
			smtp.SendEmail(u.Email, crtFile, keyFile)
		}

		//save user policy
		p.DumpJSONFile(abacPolicyFile)

		// restart apiserver to active the new PolicyFile
		_ = shell("docker restart $(docker ps | grep apiserver | awk '{print $1}')")
		//TODO: implement function by docker client:https://github.com/docker/docker/tree/master/client

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

func shell(cmd string) error {
	if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
		return err
	}
	return nil
}
