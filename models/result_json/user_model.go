package resultjson

type User struct {
	Id           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	RewardStatus string `json:"reward_status"`
	UrlUserImage string `json:"url_user_image"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	CreatedAt    string `json:"created_at"`
}
