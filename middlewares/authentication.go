package middlewares

import (
	"fmt"
	"go-gin/app/model"
	"go-gin/config"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	Id       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.StandardClaims
}

func GetUserIdFromToken(tokenReq string) (Id string) {
	token, _ := jwt.Parse(tokenReq, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	claim := token.Claims.(jwt.MapClaims)
	id := claim["id"].(string)

	return id
}

func GenerateJWTToken(user model.User) string {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	// exp := time.Now().Add(24 * time.Hour)

	expTimeMs, _ := strconv.Atoi(config.GetString("jwtExpirationMs"))
	exp := time.Now().Add(time.Millisecond * time.Duration(expTimeMs)).Unix()

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Id:       user.Id.Hex(),
		Username: user.Username,
		Roles:    user.Roles,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: exp,
			Audience:  "user",
			Issuer:    "uit",
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(config.GetString("jwtSecret")))
	if err != nil {
		logrus.Print(err)
	}
	return tokenString
}

func RequireAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		jwtToken := strings.Split(token, "Bearer ")
		// Initialize a new instance of `Claims`
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(jwtToken[1], claims, func(token *jwt.Token) (interface{}, error) {
			fmt.Println("testoor")
			fmt.Println([]byte(config.GetString("jwtSecret")))
			return []byte(config.GetString("jwtSecret")), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.Abort()
				c.Writer.WriteHeader(http.StatusUnauthorized)
			}
			c.Abort()
			c.Writer.WriteHeader(http.StatusBadRequest)
		}
		if !tkn.Valid {
			fmt.Println("not valid")
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
		}
		return
	}
}
