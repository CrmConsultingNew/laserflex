package widget

import "net/http"

func CheckWidget(w http.ResponseWriter, r *http.Request) {
	redirectURL := "https://b24app.rwp2.com/"

	// Use http.Redirect to redirect the client
	// The http.StatusFound status code is commonly used for redirects
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
