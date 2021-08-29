package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mhdiiilham/heimdall/config"
	"github.com/sirupsen/logrus"
)

var ACCESS_DENIED string = "access_denied"

type SpotifyAuthRequest struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

type SpotifyAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func SpotifyLogin(c *fiber.Ctx) error {
	return c.Redirect(fmt.Sprintf(
		"https://accounts.spotify.com/en/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s",
		config.Configuration.Spotify.ClientID,
		config.Configuration.Spotify.AuthorizationScopes,
		config.Configuration.RedirectURL,
	))
}

func SpotifyCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if state == ACCESS_DENIED {
		return c.Redirect("/")
	}

	return c.Redirect("/spotify-access-token?code=" + code)
}

func GetSpotifyAccessToken(c *fiber.Ctx) error {
	var respBody SpotifyAuthResponse
	code := c.Query("code")
	if code == "" || len(code) < 1 {
		return c.Status(http.StatusBadRequest).JSON(struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    http.StatusBadRequest,
			Message: "code cannot be empty",
		})
	}

	payload := url.Values{}
	payload.Set("grant_type", "authorization_code")
	payload.Set("code", code)
	payload.Set("redirect_uri", config.Configuration.RedirectURL)

	req, reqErr := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(payload.Encode()))
	if reqErr != nil {
		logrus.Errorf("error request: %v\n", reqErr)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.Configuration.Spotify.ClientID, config.Configuration.Spotify.ClientSecret)))))

	client := &http.Client{}
	resp, respErr := client.Do(req)
	if respErr != nil {
		logrus.Errorf("Error resp: %v\n", respErr)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&respBody)
	if err != nil {
		logrus.Error("error decoding payload: %v, err")
	}

	return c.Status(http.StatusOK).JSON(respBody)
}
