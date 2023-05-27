package rapyd

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	SaltLen         = 12
	SaltHeader      = "salt"
	TimestampHeader = "timestamp"
	AccessKeyHeader = "access_key"
	SignatureHeader = "signature"

	ContentTypeHeader  = "Content-Type"
	DefaultContentType = "application/json"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

type SignatureData struct {
	Method    string
	Path      string
	Salt      string
	Timestamp string
	AccessKey string
	SecretKey string
	Body      string
}

type Signer interface {
	signData(data SignatureData) []byte
	signRequest(r *http.Request, body []byte) error
	secretKeyString() string
	accessKeyString() string
}

type signer struct {
	accessKey []byte
	secretKey []byte
}

func NewRapydSigner(accessKey []byte, secretKey []byte) Signer {
	return &signer{
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

func (s *signer) secretKeyString() string {
	return string(s.secretKey)
}

func (s *signer) accessKeyString() string {
	return string(s.accessKey)
}

func (s *signer) signData(data SignatureData) []byte {
	if data.AccessKey == "" {
		data.AccessKey = s.accessKeyString()
	}

	if data.SecretKey == "" {
		data.SecretKey = s.secretKeyString()
	}

	toSign := data.Method + data.Path + data.Salt + data.Timestamp + data.AccessKey + data.SecretKey + data.Body

	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(toSign))
	return h.Sum(nil)
}

func (s *signer) signRequest(r *http.Request, body []byte) error {
	salt, err := getSalt(SaltLen)
	if err != nil {
		return errors.Wrap(err, "error getting salt")
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	r.Header.Add(AccessKeyHeader, s.accessKeyString())
	r.Header.Add(SaltHeader, salt)
	r.Header.Add(TimestampHeader, timestamp)
	r.Header.Add(ContentTypeHeader, DefaultContentType)

	data := SignatureData{
		Method:    strings.ToLower(r.Method),
		Path:      r.URL.RequestURI(),
		Salt:      salt,
		Timestamp: timestamp,
	}

	if len(body) != 0 {
		data.Body = string(body)
	}

	signature := s.signData(data)
	r.Header.Add(SignatureHeader, base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(signature))))
	return nil
}

func getRandomRune() (*rune, error) {
	lettersLen := len(letters)

	index, err := rand.Int(rand.Reader, new(big.Int).SetInt64(int64(lettersLen)))
	if err != nil {
		return nil, errors.Wrap(err, "error picking random number")
	}

	return &letters[index.Int64()], nil
}

// GetSalt for generating salt for certain length
func getSalt(len int) (string, error) {
	salt := make([]rune, len)

	for i := range salt {
		r, err := getRandomRune()
		if err != nil {
			return "", errors.New("error getting next rune")
		}

		salt[i] = *r
	}
	return string(salt), nil
}
