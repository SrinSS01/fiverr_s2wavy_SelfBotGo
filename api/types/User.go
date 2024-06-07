package types

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	//Discriminator        string        `json:"discriminator"`
	//PublicFlags          int           `json:"public_flags"`
	//Flags                int           `json:"flags"`
	Banner      interface{} `json:"banner"`
	AccentColor int         `json:"accent_color"`
	//GlobalName           string        `json:"global_name"`
	//AvatarDecorationData interface{}   `json:"avatar_decoration_data"`
	BannerColor string `json:"banner_color"`
	//Clan                 interface{}   `json:"clan"`
	//MfaEnabled           bool          `json:"mfa_enabled"`
	//Locale               string        `json:"locale"`
	//PremiumType          int           `json:"premium_type"`
	//Email                string        `json:"email"`
	//Verified             bool          `json:"verified"`
	//Phone                string        `json:"phone"`
	//NsfwAllowed          bool          `json:"nsfw_allowed"`
	//LinkedUsers          []interface{} `json:"linked_users"`
	//PurchasedFlags     int    `json:"purchased_flags"`
	Bio string `json:"bio"`
	//AuthenticatorTypes []int  `json:"authenticator_types"`
}
