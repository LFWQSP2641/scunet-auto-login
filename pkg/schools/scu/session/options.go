package session

type Session struct {
	PreAuth  *PreAuthInfo
	Crypto   *CryptoContext
	Cookies  map[string]string
	LoggedIn bool
}

type PreAuthInfo struct {
	LoginQuery       string
	DeviceMAC        string
	AdditionalParams map[string]string
}

type CryptoContext struct {
	PasswordPlain     string
	DeviceMAC         string
	PublicKeyExponent string
	Module            string
}
