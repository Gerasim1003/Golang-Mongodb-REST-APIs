package store

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	// "strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	// "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

//Controller ...
type Controller struct {
	Repository Repository
}

//AuthenticationMiddleware Middleware handler to handle all requests for authentication
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		autorizationHeader := r.Header.Get("authorization")

		if autorizationHeader != "" {
			bearerToken := strings.Split(autorizationHeader, " ")

			log.Println("bearer token", bearerToken)
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}

					return []byte("secret"), nil
				})

				if err != nil {
					json.NewEncoder(w).Encode(Exception{Message: err.Error()})
					return
				}

				if token.Valid {
					log.Println("TOCEN WAS VALID")
					context.Set(r, "decoded", token.Claims)
					next(w, r)
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

//GetToken  Get Authentication token GET /
func (c *Controller) GetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var user User

	json.NewDecoder(r.Body).Decode(&user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})

	log.Println("Username: ", user.Username)
	log.Println("Password: ", user.Password)

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

//Index ... GET
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	heroes := c.Repository.ReturAllHeroes(bson.M{})

	data, _ := json.Marshal(heroes)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//AddHero ...
func (c *Controller) AddHero(w http.ResponseWriter, r *http.Request) {
	var hero Hero
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	log.Println(body)

	if err != nil {
		log.Fatalln("Error AddProduct", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error AddProduct", err)
	}

	if err := json.Unmarshal(body, &hero); err != nil {
		w.WriteHeader(422)
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error AddProduct unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	log.Println(hero)
	id := c.Repository.InsertNewHero(hero)
	if id == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	return
}
