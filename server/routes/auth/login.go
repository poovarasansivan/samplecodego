package auth

import (
	"database/sql"
	"encoding/json"
	"jwt/config"
	"jwt/function"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type UserAccount struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Pass   string `json:"password"`
	Email  string `json:"email"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input UserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	var response map[string]interface{}
	var data UserAccount

	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   "Invalid Request",
		}
		function.Response(w, response)
		return
	}

	err = config.Database.QueryRow("SELECT id,username,password,email FROM login WHERE email= ? ", input.Email).Scan(&data.UserId, &data.Name, &data.Pass, &data.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			response = map[string]interface{}{
				"success": false,
				"error":   "No User Found",
			}
		} else {
			response = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		function.Response(w, response)
		return
	}

	// Check if the password is correct
	if input.Password != data.Pass {
		response = map[string]interface{}{
			"success": false,
			"error":   "Invalid Password",
		}
		function.Response(w, response)
		return
	}

	// Create JWT token
	token, err := createToken(data.UserId)
	if err != nil {
		response = map[string]interface{}{
			"success": false,
			"error":   "Error creating JWT token",
		}
		function.Response(w, response)
		return
	}

	// Include the JWT token in the response
	response = map[string]interface{}{
		"success": true,
		"user":    data,
		"token":   token,
	}
	function.Response(w, response)
}

func createToken(userID string) (string, error) {
	// Create the claims
	claims := UserClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
			Issuer:    "your-issuer", // Set your issuer
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	return token.SignedString([]byte("your-secret-key")) // Replace "your-secret-key" with your actual secret key
}

func LoginData(w http.ResponseWriter,r *http.Request){
	var response map[string]interface{}
	var users []UserAccount
	var data UserAccount
	row, err := config.Database.Query("SELECT id,username,PASSWORD,email FROM login")

	if err != nil {
		if err == sql.ErrNoRows {
			response = map[string]interface{}{
				"success": false,
				"error":   "No Request",
			}
		} else {
			response = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		function.Response(w, response)
		return
	}

	for row.Next() {
		err := row.Scan(&data.UserId, &data.Name, &data.Pass, &data.Email)
		if err != nil {
			panic(err.Error)
		}

		tempRow := UserAccount{
			UserId: data.UserId,
			Name:    data.Name,
			Pass:    data.Pass,
			Email:   data.Email,
		}
		users = append(users, tempRow)
	}
	response = map[string]interface{}{
		"success": true,
		"data":    users,
	}
	function.Response(w, response)
}