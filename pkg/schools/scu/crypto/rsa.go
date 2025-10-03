package crypto

import (
	"fmt"
	"math/big"
	S "scunet-auto-login/pkg/schools/scu/session"
)

// RSAKeyPair 表示 RSA 密钥对（只支持加密）
type RSAKeyPair struct {
	E         *big.Int
	M         *big.Int
	ChunkSize int
}

// GetKeyPair 构建 RSA 密钥对
func GetKeyPair(eHex, mHex string) (*RSAKeyPair, error) {
	e := new(big.Int)
	e.SetString(eHex, 16)

	m := new(big.Int)
	m.SetString(mHex, 16)

	chunkSize := 2 * ((m.BitLen() + 15) / 16)
	return &RSAKeyPair{
		E:         e,
		M:         m,
		ChunkSize: chunkSize,
	}, nil
}

// EncryptedString 公钥加密字符串（零填充）
func EncryptedString(key *RSAKeyPair, s string) (string, error) {
	bytes := []byte(s)
	// 填充零直到长度是 chunkSize 的倍数
	for len(bytes)%key.ChunkSize != 0 {
		bytes = append(bytes, 0)
	}

	var result []string
	for i := 0; i < len(bytes); i += key.ChunkSize {
		block := big.NewInt(0)
		for j := 0; j < key.ChunkSize && i+j < len(bytes); j += 2 {
			val := int(bytes[i+j])
			if i+j+1 < len(bytes) {
				val |= int(bytes[i+j+1]) << 8
			}
			tmp := big.NewInt(int64(val))
			tmp.Lsh(tmp, uint(16*(j/2)))
			block.Or(block, tmp)
		}

		// 使用 big.Int.Exp 直接做模幂
		crypt := new(big.Int).Exp(block, key.E, key.M)

		modulusHexLength := (key.M.BitLen() + 3) / 4
		text := fmt.Sprintf("%0*x", modulusHexLength, crypt)
		result = append(result, text)
	}

	return join(result, " "), nil
}

// EncryptedPassword 对密码进行加密
func EncryptedPassword(cryptoContext S.CryptoContext) (string, error) {
	reversed := reverseString(cryptoContext.PasswordPlain + ">" + cryptoContext.DeviceMAC)
	publicKeyExponent := cryptoContext.PublicKeyExponent
	modulus := cryptoContext.Module

	key, err := GetKeyPair(publicKeyExponent, modulus)
	if err != nil {
		return "", err
	}

	return EncryptedString(key, reversed)
}

// join 简单实现 strings.Join
func join(strs []string, sep string) string {
	res := ""
	for i, s := range strs {
		if i > 0 {
			res += sep
		}
		res += s
	}
	return res
}

// reverseString 反转字符串
func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
