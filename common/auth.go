package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AppSecret provides the necessary components for GitHub authentication as an installation.
type AppSecret struct {
	AppId          string `json:"app_id"`
	InstallationId string `json:"installation_id"`
	PrivateKey     string `json:"private_key"`
}

type accessToken struct {
	Token string `json:"token"`
}

// SecretProvider is the interface to anything that can provide the AppSecret for authentication.
type SecretProvider interface {
	Credentials() (*AppSecret, error)
}

// Authenticator is the entrypoint for authentication.
type Authenticator struct {
	secretProvider SecretProvider
}

// NewAuthenticator creates an Authenticator
func NewAuthenticator(secretProvider SecretProvider) *Authenticator {
	return &Authenticator{
		secretProvider: secretProvider,
	}
}

// ShouldRun reads stdin to determine if this authenticator should even be called. If this returns false,
// the Authenticator should not be run.
func ShouldRun() bool {
	input := make(map[string]string)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		eqIndex := strings.Index(line, "=")
		if eqIndex == -1 {
			log.Fatalf("Invalid input: %s", line)
		}

		key := strings.TrimSpace(line[:eqIndex])
		value := strings.TrimSpace(line[eqIndex+1:])
		input[key] = value
	}

	if p, ok := input["protocol"]; !ok || p != "https" {
		return false
	}

	if p, ok := input["host"]; !ok || p != "github.com" {
		return false
	}

	return true
}

// Authenticate provides the string output for a git credential command.
func (auth *Authenticator) Authenticate() (string, error) {
	creds, err := auth.secretProvider.Credentials()
	if err != nil {
		log.Fatal(err)
	}

	claims := jwt.RegisteredClaims{
		Issuer:    creds.AppId,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second * -30)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(creds.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", creds.InstallationId))
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	request := &http.Request{
		Method: "POST",
		URL:    u,
		Header: map[string][]string{
			"Authorization":        {fmt.Sprintf("Bearer %s", signed)},
			"Accept":               {"application/vnd.github+json"},
			"X-GitHub-Api-Version": {"2022-11-28"},
		},
	}

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var at accessToken
	err = json.Unmarshal(body, &at)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`protocol=https
host=github.com
capability=authtype
authtype=bearer
credential=%s
`, at.Token), nil
}
