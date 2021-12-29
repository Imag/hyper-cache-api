package types

type License struct {
	Email string `json:"email"`
	Key string `json:"key"`
	User LicenseUser `json:"user"`
	Metadata interface{} `json:"metadata"`
}

type LicenseUser struct {
	Username string `json:"username"`
	Discriminator string `json:"discriminator"`
	AvatarURL string `json:"photo_url"`
}