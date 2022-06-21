package middlewares

import "net/http"

// CORS middleware
func CORSmiddleware(w http.ResponseWriter, methods string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", methods)
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
}
