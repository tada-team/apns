package apns

import (
	"archive/zip"
	"bytes"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"path"
	"strings"

	"github.com/aai/gocrypto/pkcs7"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/crypto/pkcs12"
)

var icons = []string{
	"icon_16x16.png",
	"icon_16x16@2x.png",
	"icon_32x32.png",
	"icon_32x32@2x.png",
	"icon_128x128.png",
	"icon_128x128@2x.png",
}

type SafariOpts struct {
	Cert          []byte `json:"-"`
	IconsPath     string `json:"-"`
	AppleCertPath string `json:"-"`

	DeactivateToken func(website, token string) error `json:"-"`

	// The website name. This is the heading used in Notification Center.
	WebsiteName string `json:"websiteName"`

	// The Website push ID, as specified in your developer account.
	WebsitePushId string `json:"websitePushID"`

	// An array of websites that are allowed to request permission from the user.
	AllowedDomains []string `json:"allowedDomains"`

	//The URL to go to when the notification is clicked. Use %@ as a placeholder for
	// arguments you fill in when delivering your notification. This URL must use
	// the http or https scheme; otherwise, it is invalid.
	UrlFormatString string `json:"urlFormatString"`

	// A string that helps you identify the user. It is included in later requests
	// to your web service. This string must 16 characters or greater.
	AuthenticationToken string `json:"authenticationToken"`

	// The location used to make requests to your web service. The trailing slash should be omitted.
	WebServiceURL string `json:"webServiceURL"`
}

type handler func(w http.ResponseWriter, r *http.Request) error

func SafariUrls(opts SafariOpts) (map[string]handler, error) {
	if opts.WebServiceURL == "" {
		opts.WebServiceURL = "/push"
	}

	prefix, err := relativeUrl(opts.WebServiceURL)
	if err != nil {
		return nil, err
	}
	prefix = strings.TrimRight(prefix, "/")

	return map[string]handler{
		prefix + "/v1/log": func(w http.ResponseWriter, r *http.Request) error {
			dump, err := httputil.DumpRequest(r, true)
			if err != nil {
				return errors.Wrap(err, "dump request fail")
			}
			log.Println("apns:", r.Method, r.RequestURI, string(dump))
			io.WriteString(w, "OK")
			return nil
		},
		prefix + "/v1/pushPackages/{website}": func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Content-Disposition", "attachment; filename=\"pushPackage.zip\"")
			pushPackage, err := opts.websiteJson()
			if err != nil {
				return err
			}
			if _, err := w.Write(pushPackage); err != nil {
				return errors.Wrap(err, "send zip fail")
			}
			log.Println("apns:", r.Method, r.RequestURI)
			return nil
		},
		prefix + "/v1/devices/{device}/registrations/{website}": func(w http.ResponseWriter, r *http.Request) error {
			vars := mux.Vars(r)
			v := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if v[1] != ("ApplePushNotifications "+opts.AuthenticationToken) && v[1] != opts.AuthenticationToken {
				return fmt.Errorf("handleRegistration fail: auth: '%s'", v)
			}
			if err := opts.DeactivateToken(vars["website"], vars["device"]); err != nil {
				return errors.Wrap(err, "deactivate token fail")
			}
			io.WriteString(w, "OK")
			log.Println("apns:", r.Method, r.RequestURI)
			return nil
		},
	}, nil
}

func (opts SafariOpts) appleCert() (*x509.Certificate, error) {
	b, err := ioutil.ReadFile(opts.AppleCertPath)
	if err != nil {
		return nil, errors.Wrap(err, "invalid apple cert")
	}
	cert, err := x509.ParseCertificate(b)
	if err != nil {
		return nil, errors.Wrap(err, "x509.ParseCertificate")
	}
	return cert, nil
}

func (opts SafariOpts) websiteJson() ([]byte, error) {
	buf := new(bytes.Buffer)
	zf := zip.NewWriter(buf)
	manifest := make(map[string]string)

	for _, filename := range icons {
		data, err := ioutil.ReadFile(path.Join(opts.IconsPath, filename))
		if err != nil {
			return nil, err
		}
		if err := addToZipfile(zf, &manifest, path.Join("icon.iconset", filename), data); err != nil {
			return nil, err
		}
	}

	websiteJson, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	if err := addToZipfile(zf, &manifest, "website.json", websiteJson); err != nil {
		return nil, errors.Wrap(err, "add to zip fail")
	}

	manifestBytes, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	if err := addToZipfile(zf, nil, "manifest.json", manifestBytes); err != nil {
		return nil, errors.Wrap(err, "add to zip fail")
	}

	key, cert, err := pkcs12.Decode(opts.Cert, "")
	if err != nil {
		return nil, errors.Wrap(err, "pkcs12.Decode")
	}

	appleCert, err := opts.appleCert()
	if err != nil {
		return nil, errors.Wrap(err, "appleCert fail")
	}

	sign, err := pkcs7.Sign2(bytes.NewReader(manifestBytes), cert, key.(*rsa.PrivateKey), appleCert)
	if err != nil {
		return nil, errors.Wrap(err, "sign fail")
	}

	if err := addToZipfile(zf, nil, "signature", sign); err != nil {
		return nil, errors.Wrap(err, "add to zip fail")
	}

	if err := zf.Close(); err != nil {
		return nil, errors.Wrap(err, "zip close fail")
	}

	return buf.Bytes(), nil
}

func addToZipfile(w *zip.Writer, manifest *map[string]string, filename string, data []byte) error {
	f, err := w.Create(filename)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	if manifest != nil {
		h := sha1.New()
		h.Write(data)
		(*manifest)[filename] = hex.EncodeToString(h.Sum(nil))
	}

	return nil
}
