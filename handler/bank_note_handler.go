package handler

import (
	"banknote-tracker-auth/database"
	"banknote-tracker-auth/middlewares"
	"banknote-tracker-auth/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var (
	jwtKey = []byte("secret-key")
)

type bankNoteAuthHandler struct {
	Logger *log.Logger
}

// New creates a new banknote authentication handler
func New(logger *log.Logger) *bankNoteAuthHandler {
	return &bankNoteAuthHandler{
		Logger: logger,
	}
}

func (h *bankNoteAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	middlewares.CORSmiddleware(w, "POST, OPTION")

	loggedin, _ := IsLoggedIn(r)
	if loggedin {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "user already loggedin",
		})
		return
	}

	var data models.LoginData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "something went wrong",
		})
		return
	}

	expected, ok := database.Users[data.Username]
	if !ok || expected.Password != data.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "invalid username/password",
		})
		return
	}

	// result, err := auth.Authenticate(data)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	json.NewEncoder(w).Encode(map[string]interface{}{
	// 		"error": "invalid username/password",
	// 	})
	// 	return
	// }

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &models.Claims{
		UserNoPassword: models.UserNoPassword{
			Username: expected.Username,
			Group:    expected.Group,
			Role:     expected.Role,
		},

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Logger.Println("generated token:", tokenString)

	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
	})

	h.Logger.Println("user logged in")
}

func (h *bankNoteAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	middlewares.CORSmiddleware(w, "GET, OPTION")

	loggedIn, _ := IsLoggedIn(r)
	if !loggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "user not logged in",
		})
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "user unauthorized",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "something went wrong",
		})
		return
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	h.Logger.Println("user logged out")
	w.Write([]byte("user is logged out"))
}

func (h *bankNoteAuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	middlewares.CORSmiddleware(w, "GET, OPTION")

	loggedin, claims := IsLoggedIn(r)
	if !loggedin {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "user unauthorized",
		})
		return
	}

	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(claims.UserNoPassword)
	h.Logger.Println("user logged get user")
}

func (h *bankNoteAuthHandler) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.Login).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/logout", h.Logout).Methods(http.MethodGet)
	router.HandleFunc("/getuser", h.GetUser).Methods(http.MethodGet)
}

func IsLoggedIn(r *http.Request) (loggedIn bool, claims *models.Claims) {
	// first approach: using cookies
	cookie, err := r.Cookie("token")
	if err != nil {
		return false, nil
	}

	tokenString := cookie.Value
	if tokenString == "" {
		return false, nil
	}

	c := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return false, nil
	}

	if !token.Valid {
		return false, nil
	}

	return true, c

	// second approach: the JWT token is being sent on every request
	// var token models.Token
	// if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(map[string]interface{}{
	// 		"error": "somethin went wrong",
	// 	})
	// 	return
	// }
	// tokenString := token.Token

	// if tokenString == "" {
	// 	return false, nil
	// }

	// c := &models.Claims{}
	// token, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (interface{}, error) {
	// 	return jwtKey, nil
	// })

	// if err != nil {
	// 	return false, nil
	// }

	// if !token.Valid {
	// 	return false, nil
	// }

	// return true, c
}
