package hue

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

const (
	AuthURL  = "https://api.meethue.com/oauth2/auth"
	TokenURL = "https://api.meethue.com/oauth2/token"
	ApiURL   = "https://api.meethue.com/bridge/"
)

type Authenticator struct {
	config  *oauth2.Config
	context context.Context
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewAuthenticator(redirectURL string) Authenticator {
	appID := os.Getenv("HUE_APP_ID")
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("HUE_CLIENT_ID"),
		ClientSecret: os.Getenv("HUE_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s?appid=%s&deviceid=%s&devicename=browser", AuthURL, appID, appID),
			TokenURL: TokenURL,
		},
	}

	tr := &http.Transport{
		TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: tr})
	return Authenticator{
		config:  cfg,
		context: ctx,
	}
}

func (a *Authenticator) SetAuthInfo(clientID, secretKey string) {
	a.config.ClientID = clientID
	a.config.ClientSecret = secretKey
}

func (a *Authenticator) AuthURL(state string) string {
	return a.config.AuthCodeURL(state)
}

func (a *Authenticator) AuthURLWithOpts(state string, opts ...oauth2.AuthCodeOption) string {
	return a.config.AuthCodeURL(state, opts...)
}

func (a *Authenticator) Token(state string, r *http.Request) (*oauth2.Token, error) {
	values := r.URL.Query()
	if e := values.Get("error"); e != "" {
		return nil, errors.New("hue: auth failed - " + e)
	}
	code := values.Get("code")
	if code == "" {
		return nil, errors.New("hue: didn't get access code")
	}
	actualState := values.Get("state")
	if actualState != state {
		return nil, errors.New("hue: redirect state parameter doesn't match")
	}
	return a.config.Exchange(a.context, code)
}

func (a *Authenticator) Exchange(code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return a.config.Exchange(a.context, code, opts...)
}

func (a *Authenticator) Authenticate() (*Client, error) {
	state := fmt.Sprintf("%v", rand.Intn(10000))
	authUrl := a.AuthURL(state)
	fmt.Printf("Go to %v\n", authUrl)
	fmt.Println("Waiting 1 minute for the action")

	browser.OpenURL(authUrl)

	u, _ := url.Parse(a.config.RedirectURL)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", u.Port()),
	}

	tokenCh := make(chan *oauth2.Token)

	redirectUrl, _ := url.Parse(a.config.RedirectURL)
	http.HandleFunc(redirectUrl.Path, func(rw http.ResponseWriter, r *http.Request) {
		token, _ := a.Token(state, r)

		io.WriteString(rw, `
		<html>
			<body>
				<h1>Login successful!</h1>
				<h2>You can close this window.</h2>
			</body>
		</html>`)

		go srv.Close()

		// Hue response includes {"token_type": "BearerToken"},
		// But we need "Bearer" when you call back.
		token.TokenType = "Bearer"
		tokenCh <- token
	})

	go srv.ListenAndServe()

	var client *Client
	var token *oauth2.Token

	select {
	case token = <-tokenCh:
		client = a.NewClient(token)
		return client, nil

	case <-time.After(1 * time.Minute):
		fmt.Println("could not authenticate on time")
		go srv.Close()
		return nil, errors.New("could not authenticate on time")
	}
}

func (c *Client) Login(username string) error {
	c.clientId = username
	return nil
}

func (c *Client) CreateRemoteUser() (string, error) {
	c.logger.Info("Press the link button on the bridge")
	c.logger.Info("Waiting for 30 seconds...")

	time.Sleep(30 * time.Second)

	c.logger.Info("Enable link button...")
	err := c.EnableLinkButton()
	if err != nil {
		c.logger.Error(err, "could not enable link button")
		return "", nil
	}

	c.logger.Info("Add whitelist identifies...")
	username, err := c.AddWhitelistIdentifier()
	if err != nil {
		c.logger.Error(err, "could not add whitelist identifier")
		return "", nil
	}

	c.clientId = username

	return username, nil
}

func (c *Client) EnableLinkButton() error {
	var payload = struct {
		LinkButton bool `json:"linkedbutton"`
	}{LinkButton: true}
	req, err := c.newRequest(http.MethodPut, "0/config", payload)
	if err != nil {
		return err
	}

	resp, err := c.do(context.Background(), req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("server returned failed status")
}

func (c *Client) AddWhitelistIdentifier() (string, error) {
	var payload = struct {
		DeviceType string `json:"devicetype"`
	}{DeviceType: os.Getenv("HUE_APP_ID")}
	req, err := c.newRequest(http.MethodPost, "", payload)
	if err != nil {
		return "", err
	}

	apiResponses := new([]ApiResponse)
	resp, err := c.do(context.Background(), req, apiResponses)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		if (*apiResponses)[0].Success != nil {
			return (*apiResponses)[0].Success["username"].(string), nil
		} else {
			return "", errors.New((*apiResponses)[0].Error.Description)
		}
	}

	return "", errors.New("server returned failed status")
}

func (a *Authenticator) NewClient(token *oauth2.Token) *Client {
	httpClient := a.config.Client(a.context, token)
	client, err := newClient(ApiURL, &ClientOptions{HttpClient: httpClient})
	if err != nil {
		return nil
	}
	return client
}
