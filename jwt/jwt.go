package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("hello-jwt")

func main() {
	token1 := CreateToken("hello-user")
	fmt.Println("token > ", token1)

	claims1 := getClaimFromToken(token1)
	fmt.Println("claims1 > ", claims1)

	claim2 := getClaimFromToken("eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJhZG1pbiIsImluZm9zIjp7ImVudE5hbWUiOiLshJzsmrjtirnrs4Tsi5wiLCJ1bml0TmFtZSI6IuygleuztO2ZlOu2gOyEnCIsImlwIjoiMTI3LjAuMC4xIiwiY29udHJhY3QiOiIwMTAtMTIzNC05OTk5IiwiYWN0aXZlIjoibG9jYWwiLCJpc1NTT0xvZ2luIjpmYWxzZSwibG9naW5Vc2VyTmFtZSI6Iuq0gOumrOyekCIsImlzT1NNYW5hZ2VyIjp0cnVlLCJsb2dpblRpbWUiOjE2NjM3NTI4NTY5NzksIm9wZW5zdGFja19jcmVkZW50aWFsIjoiOWI5MjYzN2MtY2I5Yi00NDMyLTgwMDYtYTZkNDI0MTQzNzEyIiwiZW50Q29kZSI6IlNFTCIsInVuaXRDb2RlIjoiMDAyIiwiZW1haWwiOiJjcm1vb25AaW5ub2dyaWQuY29tIiwidm1fYXBwbGljYXRpb25fY2xvdWRfdHlwZSI6Ik9QRU5TVEFDSyIsInVzZXJuYW1lIjoiYWRtaW4ifSwidXNlcklkIjoic2VjbG91ZGl0LWFkbWluIiwiaWF0IjoxNjYzNzUyODU3LCJleHAiOjE2OTUyODg4NTd9.Q8Vy7soUrM_rXBzNOnC3OkLBSY1bw5K8gMJChmDQKjM")
	fmt.Println(claim2)
}

func CreateToken(userId string) (tokenString string) {
	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenClaims.SignedString(jwtKey)

	if err != nil {
		fmt.Println("failed create jwt token", err)
	}
	return tokenString
}

func getClaimFromToken(tokenString string) (claims *Claims) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil
	}
	return claims
}

func GetUserId(tokenString string) string {
	claims := getClaimFromToken(tokenString)
	return claims.UserID
}
