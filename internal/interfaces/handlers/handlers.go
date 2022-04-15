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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token not found", http.StatusBadRequest)
		}
		splitToken := strings.Split(authHeader, "Bearer ")
		authToken := splitToken[1]

		token, err := jwt.ParseWithClaims(authToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(verifyKey), nil
		})

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			logger.Info.Printf("Get request from account with id %v. Token expire in %v", claims.UserId, claims.StandardClaims.ExpiresAt.Unix())
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
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

func getBriefcase(app storage.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		logger.Info.Println("Got new request for getting briefcase")
		errorMessage := "Error getting briefcase"

		profileId := context.Get(r, "profile_id").(int)
		briefcase, err := app.GetBriefcase(profileId)
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
		resp["data"] = briefcase
		json.NewEncoder(w).Encode(resp)
	})
}

func payBalance(app storage.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		logger.Info.Println("Got new request for adding money to balance")
		errorMessage := "Error updating balance"

		profileId := context.Get(r, "profile_id").(int)
		var requestBody UpdateBalanceBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			logger.Error.Printf("Failed to decode balance body. Got %v", err)
			http.Error(w, "wrong structure", http.StatusBadRequest)
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

func buyStocks(app storage.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		logger.Info.Println("Got new request for buying stocks")

		profileId := context.Get(r, "profile_id").(int)
		var requestBody StocksBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			logger.Error.Printf("Failed to decode stocks. Got %v", err)
			http.Error(w, "wrong structure", http.StatusBadRequest)
			return
		}

		err = app.BuyStocks(profileId, requestBody.Ticker, requestBody.Amount)
		switch err {
		case storage.AccountNotFoundError, storage.NotEnoughMoneyError:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case storage.InternalError:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		resp := make(map[string]interface{})
		resp["status"] = "success"
		json.NewEncoder(w).Encode(resp)
	})
}

func sellStocks(app storage.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		logger.Info.Println("Got new request for buying stocks")

		profileId := context.Get(r, "profile_id").(int)
		var requestBody StocksBody

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			logger.Error.Printf("Failed to decode stocks. Got %v", err)
			http.Error(w, "wrong structure", http.StatusBadRequest)
			return
		}

		err = app.SellStocks(profileId, requestBody.Ticker, requestBody.Amount)
		switch err {
		case storage.AccountNotFoundError, storage.NotEnoughStocksError:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case storage.InternalError:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		resp := make(map[string]interface{})
		resp["status"] = "success"
		json.NewEncoder(w).Encode(resp)
	})
}

func Make(r *mux.Router, app storage.Controller) {
	apiUri := "/api"
	r.Use(jwtHandler)
	serviceRouter := r.PathPrefix(apiUri).Subrouter()
	serviceRouter.Handle("/account/balance", getBalance(app)).Methods("GET")
	serviceRouter.Handle("/account/balance", payBalance(app)).Methods("POST")
	serviceRouter.Handle("/briefcase", getBriefcase(app)).Methods("GET")
	serviceRouter.Handle("/briefcase/buy", buyStocks(app)).Methods("POST")
	serviceRouter.Handle("/briefcase/sell", sellStocks(app)).Methods("POST")
}
