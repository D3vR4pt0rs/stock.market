package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/D3vR4pt0rs/logger"
	"github.com/dgrijalva/jwt-go/v4"
	"market/internal/usecases/storage"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func jwtHandler(next http.Handler) http.Handler {
	verifyKey := os.Getenv("SECRET_KEY")
	logger.Info.Printf(verifyKey)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splitToken := strings.Split(authHeader, "Bearer ")
		authToken := splitToken[1]

		token, err := jwt.ParseWithClaims(authToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(verifyKey), nil
		})

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			logger.Info.Printf("%v %v", claims.UserId, claims.StandardClaims.ExpiresAt.Unix())
			context.Set(r, "profile_id", claims.UserId)
			next.ServeHTTP(w, r)
		} else {
			logger.Error.Println(err)
			http.Error(w, "Token is wrong", http.StatusUnauthorized)
		}
	})
}

func getBalance(app storage.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Got new request for getting balance")
		errorMessage := "Error getting balance"

		profileId := context.Get(r, "profile_id").(int)
		balance, err := app.GetBalance(profileId)
		switch err {
		case storage.AccountNotFoundError:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case storage.InternalError:
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		resp := make(map[string]interface{})
		resp["balance"] = balance
		json.NewEncoder(w).Encode(resp)
	})
}

func payBalance(app storage.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info.Println("Got new request for adding money to balance")
		errorMessage := "Error updating balance"

		profileId := context.Get(r, "profile_id").(int)
		var requestBody UpdateBalanceBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			logger.Error.Printf("Failed to decode credentials. Got %v", err)
			http.Error(w, errorMessage, http.StatusBadRequest)
			return
		}

		balance, err := app.UpdateBalance(profileId, requestBody.Amount)
		switch err {
		case storage.AccountNotFoundError:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case storage.InternalError:
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		resp := make(map[string]interface{})
		resp["balance"] = balance
		json.NewEncoder(w).Encode(resp)
	})
}

func Make(r *mux.Router, app storage.Controller) {
	apiUri := "/api"
	r.Use(jwtHandler)
	serviceRouter := r.PathPrefix(apiUri).Subrouter()
	serviceRouter.Handle("/account/balance", getBalance(app)).Methods("GET")
	serviceRouter.Handle("/account/balance", payBalance(app)).Methods("POST")
}
