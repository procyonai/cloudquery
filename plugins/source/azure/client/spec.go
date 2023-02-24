package client

type Creds struct {
	TenantID     string `json:"tenantID"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

type Spec struct {
	Subscriptions []string `json:"subscriptions"`
	Credentials   *Creds   `json:"credentials"`
}
