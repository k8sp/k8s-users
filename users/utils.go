package users

// Users describe http request body
// url:     http://<domain>/users
// method:  POST
type Users struct {
	Username  string
	Namespace string
	Email     string
}

func getDefaultServiceAcccount(namespace string) string {
	return "system:serviceaccount:" + namespace + ":default"
}
