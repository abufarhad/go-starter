package middlewares

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/abufarhad/golang-starter-rest-api/domain/dto"
	"github.com/abufarhad/golang-starter-rest-api/internal/logger"
	"github.com/abufarhad/golang-starter-rest-api/internal/utils/methodutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(*gin.Context) bool

	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFunc func(*gin.Context)

	// JWTSuccessHandler defines a function which is executed for a valid token.
	JWTSuccessHandler func(*gin.Context)

	// JWTErrorHandler defines a function which is executed for an invalid token.
	JWTErrorHandler func(error) error

	// JWTErrorHandlerWithContext is almost identical to JWTErrorHandler, but it's passed the current context.
	JWTErrorHandlerWithContext func(error, *gin.Context) error

	jwtExtractor func(*gin.Context) (string, error)

	// JWTConfig defines the config for JWT middleware.
	JWTConfig struct {
		Skipper                 Skipper
		BeforeFunc              BeforeFunc
		SuccessHandler          JWTSuccessHandler
		ErrorHandler            JWTErrorHandler
		ErrorHandlerWithContext JWTErrorHandlerWithContext
		SigningKey              interface{}
		SigningKeys             map[string]interface{}
		SigningMethod           string
		ContextKey              string
		Claims                  jwt.Claims
		TokenLookup             string
		AuthScheme              string
		keyFunc                 jwt.Keyfunc
	}
)

// Algorithms
const (
	AlgorithmHS256 = "HS256"
)

// Errors
var (
	ErrJWTMissing = errors.New("missing, invalid, or expired jwt")
)

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(c *gin.Context) bool {
	return false
}

// DefaultJWTConfig is the default JWT auth middleware config.
var DefaultJWTConfig = JWTConfig{
	Skipper:       DefaultSkipper,
	SigningMethod: AlgorithmHS256,
	ContextKey:    "auth",
	TokenLookup:   "header:" + "Authorization",
	AuthScheme:    "Bearer",
	Claims:        jwt.MapClaims{},
}

// JWT returns a JSON Web Token (JWT) auth middleware.
func JWT(key interface{}) gin.HandlerFunc {
	c := DefaultJWTConfig
	c.SigningKey = key
	return JWTWithConfig(c)
}

// JWTWithConfig returns a JWT auth middleware with config.
func JWTWithConfig(config JWTConfig) gin.HandlerFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultJWTConfig.Skipper
	}
	if config.SigningKey == nil && len(config.SigningKeys) == 0 {
		panic("gin: jwt middleware requires signing key")
	}
	if config.SigningMethod == "" {
		config.SigningMethod = DefaultJWTConfig.SigningMethod
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultJWTConfig.ContextKey
	}
	if config.Claims == nil {
		config.Claims = DefaultJWTConfig.Claims
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultJWTConfig.TokenLookup
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultJWTConfig.AuthScheme
	}
	config.keyFunc = func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != config.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		if len(config.SigningKeys) > 0 {
			if kid, ok := t.Header["kid"].(string); ok {
				if key, ok := config.SigningKeys[kid]; ok {
					return key, nil
				}
			}
			return nil, fmt.Errorf("unexpected jwt key id=%v", t.Header["kid"])
		}

		return config.SigningKey, nil
	}

	parts := strings.Split(config.TokenLookup, ":")
	extractor := jwtFromHeader(parts[1], config.AuthScheme)
	switch parts[0] {
	case "query":
		extractor = jwtFromQuery(parts[1])
	case "param":
		extractor = jwtFromParam(parts[1])
	case "cookie":
		extractor = jwtFromCookie(parts[1])
	}

	return func(c *gin.Context) {
		if config.Skipper(c) {
			c.Next()
			return
		}

		auth, err := extractor(c)
		if err != nil {
			if config.ErrorHandler != nil {
				config.ErrorHandler(err)
				return
			}
			if config.ErrorHandlerWithContext != nil {
				config.ErrorHandlerWithContext(err, c)
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing, invalid, or expired jwt"})
			c.Abort()
			return
		}

		var token *jwt.Token
		if _, ok := config.Claims.(jwt.MapClaims); ok {
			token, err = jwt.Parse(auth, config.keyFunc)
		} else {
			t := reflect.ValueOf(config.Claims).Type().Elem()
			claims := reflect.New(t).Interface().(jwt.Claims)
			token, err = jwt.ParseWithClaims(auth, claims, config.keyFunc)
		}

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrJWTMissing})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrJWTMissing})
			c.Abort()
			return
		}

		tokenDetails := &dto.JwtToken{}
		if err := methodutil.MapToStruct(claims, tokenDetails); err != nil {
			logger.Error(err.Error(), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrJWTMissing})
			c.Abort()
			return
		}

		c.Set(config.ContextKey, &dto.LoggedInUser{
			ID:          int(tokenDetails.UserID),
			AccessUuid:  tokenDetails.AccessUuid,
			RefreshUuid: tokenDetails.RefreshUuid,
		})

		c.Next()
	}
}

func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c *gin.Context) (string, error) {
		auth := c.Request.Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}

func jwtFromQuery(param string) jwtExtractor {
	return func(c *gin.Context) (string, error) {
		token := c.Query(param)
		if token == "" {
			return "", ErrJWTMissing
		}
		return token, nil
	}
}

func jwtFromParam(param string) jwtExtractor {
	return func(c *gin.Context) (string, error) {
		token := c.Param(param)
		if token == "" {
			return "", ErrJWTMissing
		}
		return token, nil
	}
}

func jwtFromCookie(name string) jwtExtractor {
	return func(c *gin.Context) (string, error) {
		cookie, err := c.Cookie(name)
		if err != nil {
			return "", ErrJWTMissing
		}
		return cookie, nil
	}
}
