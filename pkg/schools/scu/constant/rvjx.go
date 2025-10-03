package constant

const RSAPublicKeyExponent = "10001"
const RSAModulus = "94dd2a8675fb779e6b9f7103698634cd400f27a154afa67af6166a43fc26417222a79506d34cacc7641946abda1785b7acf9910ad6a0978c91ec84d40b71d2891379af19ffb333e7517e390bd26ac312fe940c340466b4a5d4af1d65c3b5944078f96a1a51a5a53e4bc302818b7c9f63c4a1b07bd7d874cef1c3d4b2f5eb7871"

const MainUrl = "http://192.168.2.135/"
const LoginPostUrl = MainUrl + "eportal/InterFace.do?method=login"

const (
	Accept         = "*/*"
	AcceptEncoding = "gzip, deflate"
	ContentType    = "application/x-www-form-urlencoded; charset=UTF-8"
	UserAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"
)

func HttpHeader() map[string]string {
	return map[string]string{
		"Accept":          Accept,
		"Accept-Encoding": AcceptEncoding,
		"Content-Type":    ContentType,
		"User-Agent":      UserAgent,
	}
}

// 服务类型常量
const (
	ServiceChinaTelecom = "%E7%94%B5%E4%BF%A1%E5%87%BA%E5%8F%A3" // 电信出口
	ServiceChinaMobile  = "%E7%A7%BB%E5%8A%A8%E5%87%BA%E5%8F%A3" // 移动出口
	ServiceChinaUnicom  = "%E8%81%94%E9%80%9A%E5%87%BA%E5%8F%A3" // 联通出口
	ServiceEduNet       = "internet"
)

func ServiceCodes() map[string]string {
	return map[string]string{
		"CHINATELECOM": ServiceChinaTelecom,
		"CHINAMOBILE":  ServiceChinaMobile,
		"CHINAUNICOM":  ServiceChinaUnicom,
		"EDUNET":       ServiceEduNet,
	}
}
