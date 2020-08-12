package apns

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/pkcs12"

	"golang.org/x/net/http2"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type Config struct {
	Host         string
	Bundle       string
	KeyId        string
	TeamId       string
	AuthKey      []byte
	SafariCert   []byte
	mux          sync.Mutex
	authKey      *ecdsa.PrivateKey
	tokenValue   *string
	generated    *time.Time
	safariClient *http.Client
}

const urlMask = "https://%s/3/device/%s"

func (c *Config) Send(url string, req interface{}, headers Headers, client *http.Client) (r Result) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		r.Code = FailNow
		r.Error = errors.Errorf("json fail on: %v", req)
		return
	}

	maxPayloadSize := 4096
	if headers.PushType == PushTypeVoip {
		maxPayloadSize = 5120
	}

	if len(reqBytes) > maxPayloadSize {
		r.Code = FailNow
		r.Error = errors.Errorf("big json: %v", string(reqBytes))
		return
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		r.Code = RetryNow
		r.Error = errors.Wrap(err, "post fail")
		return
	}

	token, err := c.getToken()
	if err != nil {
		r.Code = FailNow
		r.Error = errors.Wrap(err, "get token fial")
		return
	}

	headers.topic = c.Bundle
	for k, v := range headers.Map() {
		request.Header.Set(k, v)
	}

	request.Header.Set("Authorization", "bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	if client == nil {
		client = http.DefaultClient
	}
	response, err := client.Do(request)
	if err != nil {
		r.Code = RetryNow
		r.Error = errors.Wrap(err, "client do fail")
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		r.Code = RetryNow
		r.Error = errors.Wrap(err, "read all fail")
		return
	}

	r.DebugRequest = fmt.Sprintf("url: %s\n%s", url, string(reqBytes))
	r.DebugResponse = fmt.Sprintf("code: %d\nbody: %s", response.StatusCode, string(body))

	if response.StatusCode == 200 {
		r.Code = Ok
		return
	}

	if 400 <= response.StatusCode && response.StatusCode <= 499 {
		data := new(struct {
			Reason string `json:"reason"`
		})

		if err := json.Unmarshal(body, data); err != nil {
			r.Code = RetryLater
			r.Error = errors.Wrap(err, "json error")
			return
		}

		switch data.Reason {
		case "BadDeviceToken", "Unregistered", "TopicDisallowed", "DeviceTokenNotForTopic", "InvalidProviderToken":
			r.Code = InvalidConfig
		case "ExpiredProviderToken":
			c.resetToken()
			r.Code = RetryNow
		case "TooManyRequests":
			r.Code = RetryLater
		case "TooManyProviderTokenUpdates":
			c.resetToken()
			r.Code = RetryLater
		case "MissingProviderToken":
			r.Code = FailNow
		default:
			r.Code = FailNow
		}

		r.Error = fmt.Errorf(data.Reason)
		return
	}

	if 500 <= response.StatusCode && response.StatusCode <= 599 {
		r.Code = RetryLater
		r.Error = fmt.Errorf("5xx")
		return
	}

	r.Code = FailNow
	r.Error = fmt.Errorf("unknown response")
	return
}

func (c *Config) resetToken() {
	log.Println("apns: reset token")
	c.generated = nil
}

func (c *Config) getToken() (string, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.generated == nil || time.Since(*c.generated) > 59*time.Minute {
		ts := time.Now()
		token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
			"iss": c.TeamId,
			"iat": ts.Unix(),
		})
		token.Header["kid"] = c.KeyId

		key, err := c.getAuthKey()
		if err != nil {
			return "", err
		}

		val, err := token.SignedString(key)
		if err != nil {
			return "", errors.Wrap(err, "token signing fail")
		}

		c.tokenValue = &val
		c.generated = &ts
	}

	return *c.tokenValue, nil
}

func (c *Config) getAuthKey() (*ecdsa.PrivateKey, error) {
	if c.authKey == nil {
		block, _ := pem.Decode(c.AuthKey)
		if block == nil {
			return nil, fmt.Errorf("pem decode error")
		}
		pKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.Wrap(err, "p8 parse error")
		}
		pkey, ok := pKey.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not ECDSA private key")
		}
		c.authKey = pkey
	}
	return c.authKey, nil
}

func (c *Config) getSafariClient() (*http.Client, error) {
	if c.safariClient == nil {
		// https://support.airship.com/hc/en-us/articles/360017992631-How-to-make-an-Apple-Safari-Web-Push-certificate
		key, cert, err := pkcs12.Decode(c.SafariCert, "")
		if err != nil {
			return nil, errors.Wrap(err, "pkcs12.Decode")
		}
		tlsCert := tls.Certificate{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  key.(*rsa.PrivateKey),
			Leaf:        cert,
		}
		tlsCfg := &tls.Config{Certificates: []tls.Certificate{tlsCert}}
		tlsCfg.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsCfg}
		if err := http2.ConfigureTransport(transport); err != nil {
			return nil, errors.Wrap(err, "ConfigureTransport fail")
		}
		c.safariClient = &http.Client{Transport: transport}
	}
	return c.safariClient, nil
}
