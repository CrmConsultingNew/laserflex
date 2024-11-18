package endpoints

import (
	"fmt"
	"log"
	"net/http"
)

const (
	RedirectURI = "https://onviz-api.ru/redir"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf("OK")))
	if err != nil {
		log.Println("Error redirect", err.Error())
		return
	}

	fmt.Println(r.Body)
}
