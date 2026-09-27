package inputjson

type UserLogin struct {
	EmailOrPhone string `json:"emailorphone" form:"emailorphone"`
	Password     string `json:"password" form:"password"`
}
