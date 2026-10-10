package models

type PasswordCredentials struct {
	Password string `json:"password" form:"password"`
}

type User struct {
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Email     string `json:"email" form:"email"`
	Phone     string `json:"phone" form:"phone"`
	PasswordCredentials
}

type UserLogin struct {
	EmailOrPhone string `json:"emailorphone" form:"emailorphone"`
	PasswordCredentials
}

type UserNoPassword struct {
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Email     string `json:"email" form:"email"`
	Phone     string `json:"phone" form:"phone"`
}

type VerifyFirebaseToken struct {
	FirebaseToken string `json:"firebase_token" form:"token"`
}
