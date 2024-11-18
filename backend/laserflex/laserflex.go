package laserflex

import (
	"io"
	"log"
	"net/http"
)

func LaserflexGetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection is starting...")
	// Чтение тела запроса
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	log.Println("FULL REQUEST CONNECTION: ", string(bs))
}
