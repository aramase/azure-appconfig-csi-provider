package appconfig

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// This file was taken from https://github.com/Azure/appconfig-go/blob/56a54ac4dafeabdc36460b59617c89bb71bcae3b/sign.go
// and modified to fit the needs of this project.

//SignRequest Setup the auth header for accessing Azure AppConfiguration service
func SignRequest(req *http.Request, id, secret string) error {
	method := req.Method
	host := req.URL.Host
	pathAndQuery := req.URL.Path
	if req.URL.RawQuery != "" {
		pathAndQuery = pathAndQuery + "?" + req.URL.RawQuery
	}

	content := []byte{}
	var err error
	if req.Body != nil {
		content, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
	}

	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return err
	}

	timestamp := time.Now().UTC().Format(http.TimeFormat)
	contentHash := getContentHashBase64(content)
	stringToSign := fmt.Sprintf("%s\n%s\n%s;%s;%s", strings.ToUpper(method), pathAndQuery, timestamp, host, contentHash)
	signature := getHmac(stringToSign, key)

	req.Header.Set("x-ms-content-sha256", contentHash)
	req.Header.Set("x-ms-date", timestamp)
	req.Header.Set("Authorization", "HMAC-SHA256 Credential="+id+", SignedHeaders=x-ms-date;host;x-ms-content-sha256, Signature="+signature)

	return nil
}

func getContentHashBase64(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func getHmac(content string, key []byte) string {
	hmac := hmac.New(sha256.New, key)
	hmac.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}
