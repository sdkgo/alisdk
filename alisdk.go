package alisdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

func PopSignature(params map[string]string, secret string) (string, string) {
	keys := make([]string, len(params))
	i := 0
	for key := range params {
		keys[i] = key
		i++
	}
	sort.Strings(keys[:])
	buf := new(bytes.Buffer)
	for i, key := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(specialUrlEncode(key) + "=" + specialUrlEncode(params[key]))
	}
	paramStr := buf.String()
	forSignatureString := "GET&%2F&" + specialUrlEncode(paramStr)

	mac := hmac.New(sha1.New, []byte(secret+"&"))
	mac.Write([]byte(forSignatureString))
	signBytes := mac.Sum(nil)
	signedStr := base64.StdEncoding.EncodeToString(signBytes)
	return signedStr, paramStr
}

var replacer = strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")

func specialUrlEncode(value string) string {
	result := url.QueryEscape(value)
	return replacer.Replace(result)
}

func randomUuid() string {
	var u [16]byte
	rand.Reader.Read(u[:])
	//u[6] = (u[6] & 0x0f) | (4 << 4)
	//u[8] = u[8]&(0xff>>2) | (0x02 << 6)
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}
