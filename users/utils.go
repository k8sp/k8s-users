package users

// UsersBody describe http request body
// url:     http://<domain>/users
// method:  POST
type UsersBody struct {
	Username  string
	Namespace string
	Email     string
}

func getDefaultServiceAcccount(username string) string {
	return "system:serviceaccount:" + username + ":default"
}
