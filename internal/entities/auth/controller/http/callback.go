package http

import (
	"errors"
	store2 "github.com/Base-111/backend/internal/entities/auth/store"
	"github.com/Base-111/backend/pkg/errors/api"
	"log/slog"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func (a *AuthHandler) Callback(c *gin.Context) {
	callbackData, err := newCallbackData(c)
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})
		return
	}
	if err = callbackData.verify(c, a.authStore); err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})
		return
	}
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("grant_type", "authorization_code"),
	}
	var oauth2Token *oauth2.Token
	oauth2Token, err = a.authClient.Oauth.Exchange(c, callbackData.authzCode, opts...)
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})
		return
	}
	var oidcToken *oidcToken
	oidcToken, err = newOIDCToken(oauth2Token, a.authClient.OIDC)
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})
		return
	}
	var userInfoClaims *userInfoClaims
	userInfoClaims, err = oidcToken.getClaims(c)
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})
		return
	}
	sessionData := store2.SessionData{
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		UserInfoData: &store2.UserInfoData{
			Email:    userInfoClaims.Email,
			FullName: userInfoClaims.Name,
		},
	}
	if err = a.sessionStore.SaveSession(c, userInfoClaims.Sub, &sessionData); err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})
		return
	}

	c.SetCookie(
		"session_id",       // name
		userInfoClaims.Sub, // value (user ID)
		3600,               // maxAge (1 hour)
		"/",                // path
		"",                 // domain
		true,               // secure
		true,               // httpOnly
	)

	c.SetCookie(
		"user_email",
		userInfoClaims.Email,
		3600,
		"/",
		"",
		true,
		true,
	)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

type callbackData struct {
	stateID   string
	authzCode string
}

func newCallbackData(c *gin.Context) (*callbackData, error) {
	stateID := c.Query("state")
	if stateID == "" {
		slog.WarnContext(c, "stateID is required")
		return nil, errors.New("stateID is required")
	}
	authorizationCode := c.Query("code")
	if authorizationCode == "" {
		slog.WarnContext(c, "authorizationCode is required")
		return nil, errors.New("authorizationCode is required")
	}
	return &callbackData{
		stateID:   stateID,
		authzCode: authorizationCode,
	}, nil
}
func (c *callbackData) verify(ctx *gin.Context, authStore store2.AuthStore) error {
	stateIDData, err := authStore.GetState(ctx, c.stateID)
	if err != nil {
		return err
	}
	if stateIDData != c.stateID {
		slog.ErrorContext(ctx, "invalid stateID")
		return errors.New("invalid stateID")
	}
	if err = authStore.DeleteState(ctx, stateIDData); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return err
	}
	return nil
}

type oidcToken struct {
	rawIDToken string
	verifier   *oidc.IDTokenVerifier
}

func newOIDCToken(oauthToken *oauth2.Token, verifier *oidc.IDTokenVerifier) (*oidcToken, error) {
	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		slog.Error("no id_token in response")
		return nil, errors.New("no id_token in response")
	}
	return &oidcToken{
		rawIDToken: rawIDToken,
		verifier:   verifier,
	}, nil
}

type userInfoClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Sub   string `json:"sub"`
}

func (o *oidcToken) getClaims(c *gin.Context) (*userInfoClaims, error) {
	idToken, err := o.verifier.Verify(c, o.rawIDToken)
	if err != nil {
		slog.Error("failed to verify ID token")
		return nil, errors.New("failed to verify ID token")
	}
	userInfoClaims := &userInfoClaims{}
	if err := idToken.Claims(&userInfoClaims); err != nil {
		slog.Error("failed to extract claims")
		return nil, errors.New("failed to extract claims")
	}
	return userInfoClaims, nil
}
