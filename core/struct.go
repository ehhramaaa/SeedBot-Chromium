package core

import (
	"net"
	"net/http"

	"github.com/go-rod/rod"
)

type Client struct {
	Account     Account
	Proxy       string
	AccessToken string
	HttpClient  *http.Client
	Browser     *rod.Browser
}

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type Account struct {
	Phone          string
	QueryId        string
	UserId         int
	Username       string
	FirstName      string
	LastName       string
	AuthDate       string
	Hash           string
	AllowWriteToPm bool
	LanguageCode   string
	QueryData      string
}
