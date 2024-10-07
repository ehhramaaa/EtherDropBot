package core

import (
	"net"
	"net/http"

	"github.com/go-rod/rod"
)

type Client struct {
	account     Account
	proxy       string
	accessToken string
	httpClient  *http.Client
	browser     *rod.Browser
}

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type Account struct {
	phone          string
	queryId        string
	userId         int
	username       string
	firstName      string
	lastName       string
	authDate       string
	hash           string
	allowWriteToPm bool
	languageCode   string
	queryData      string
	walletAddress  string
}
