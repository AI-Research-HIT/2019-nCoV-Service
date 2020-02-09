package auth

import (
	"time"

	"github.com/adam-hanna/jwt-auth/jwt"
	"github.com/ender-wan/ewlog"
)

var Jwt jwt.Auth

func InitJwt() {
	err := jwt.New(&Jwt, jwt.Options{
		SigningMethodString:   "RS256",
		PrivateKeyLocation:    "2019-nCoV-Service-oauth.rsa",
		PublicKeyLocation:     "2019-nCoV-Service-oauth.rsa.pub",
		RefreshTokenValidTime: 72 * time.Hour,
		AuthTokenValidTime:    72 * time.Hour,
		Debug:                 false,
		IsDevEnv:              true,
		BearerTokens:          false,
	})
	if err != nil {
		ewlog.Fatal(err)
	}
}

// func IssueToken(w http.ResponseWriter, user *protodef.User) (err error) {
// 	jsonBytes, err := json.Marshal(&user)
// 	if err != nil {
// 		return
// 	}

// 	var customClaims map[string]interface{}
// 	err = json.Unmarshal(jsonBytes, &customClaims)
// 	if err != nil {
// 		return
// 	}

// 	claims := jwt.ClaimsType{}
// 	claims.StandardClaims.Id = user.Id
// 	claims.CustomClaims = customClaims

// 	err = Jwt.IssueNewTokens(w, &claims)

// 	return
// }
