package token

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type OpenIdConfig struct {
	JwksUri string `json:"jwks_uri"`
}

type Keys struct {
	Keys []Key `json:"keys"`
}

type Key struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5C []string `json:"x5c"`
}

type keyCache struct {
	counter uint
	value   string
}

const callsBeforeExpiringCache uint = 10000

var getPEMPublicKey = getPEMPublicKeyCacheAware()

func Validate(jwtToken string) (*jwt.Token, error) {
	tkn, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid := token.Header["kid"].(string)

		cert, err := getPEMPublicKey(kid)
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, err
		}

		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return tkn, nil
}

func getPEMPublicKeyCacheAware() func(kid string) (string, error) {
	var cache = make(map[string]keyCache)
	return func(kid string) (string, error) {

		if ret, ok := cache[kid]; ok && ret.counter < callsBeforeExpiringCache {
			fmt.Println("Getting public key from cache")

			cache[kid] = keyCache{counter: ret.counter + 1, value: ret.value}
			return ret.value, nil
		}
		fmt.Println("Getting public key live")
		const kAzureOpenIDConfiguration = "https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration"

		// get OpenId configuration
		resp, err := http.Get(kAzureOpenIDConfiguration)
		if err != nil {
			return "", err
		}

		//read the body response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var openidConfig OpenIdConfig
		//extract the url containing the public keys

		if err = json.Unmarshal(body, &openidConfig); err != nil {
			return "", err
		}

		// http request to the keys url
		resp, err = http.Get(openidConfig.JwksUri)
		if err != nil {
			return "", err
		}

		//read the body response
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var keys Keys

		if err = json.Unmarshal(body, &keys); err != nil {
			return "", err
		}

		var publicKey Key
		for _, key := range keys.Keys {
			if key.Kid == kid {
				publicKey = key
				break
			}
		}

		//create the PEM file
		certificate := "-----BEGIN PUBLIC KEY-----\n" + publicKey.X5C[0] + "\n-----END PUBLIC KEY-----"
		cache[kid] = keyCache{value: certificate}
		return certificate, nil
	}
}
