package profile

type Profile struct {
	Domain   string `json:"domain"`
	Port     uint16 `json:"port"`
	Insecure bool   `json:"insecure"`
}
