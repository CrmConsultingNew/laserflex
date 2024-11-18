package authorize

/*
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	clientID     = "local.660318677effc3.11656865"
	redirectURI  = "https://www.applicationhost.com/application/"
	oauthBaseURL = "https://oauth.bitrix.info"
)

func startOAuthFlow(w http.ResponseWriter, r *http.Request) {
	// Generate the Bitrix24 authorization URL
	authURL := fmt.Sprintf("%s/oauth/authorize/?client_id=%s&redirect_uri=%s&state=%s",
		oauthBaseURL, clientID, redirectURI, "yourStateValue") // Customize yourStateValue as needed

	// Redirect the user to Bitrix24's authorization page
	http.Redirect(w, r, authURL, http.StatusFound)
}

func handleAuthRedirect(w http.ResponseWriter, r *http.Request) {
	// Extract the code and state from query parameters
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Optionally, verify the state parameter to protect against CSRF attacks

	// Exchange the authorization code for an access token
	token, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Access Token: %s", token)
}

func exchangeCodeForToken(code string) (string, error) {
	// Prepare the request to exchange the authorization code for an access token
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", "LJSl0lNB76B5YY6u0YVQ3AW0DrVADcRTwVr4y99PXU1BWQybWK")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	tokenURL := fmt.Sprintf("%s/oauth/token/", oauthBaseURL)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}
*/
