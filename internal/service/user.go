package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/Xapsiel/bpla_dashboard/internal/config"
	"github.com/Xapsiel/bpla_dashboard/internal/model"
)

type UserService struct {
	repo     Repository
	config   *oauth2.Config
	OIDC     *oidc.IDTokenVerifier
	provider *oidc.Provider
}

func NewUserService(repo Repository, cfg config.OidcConfig) *UserService {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Timeout: 10 * time.Second, Transport: tr}

	ctx := oidc.ClientContext(context.Background(), client)
	provider, err := oidc.NewProvider(ctx, fmt.Sprintf("%s/realms/%s", cfg.KeycloakURL, cfg.KeycloakRealm))
	if err != nil {
		return nil
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.KeycloakSecret,
		RedirectURL:  cfg.RedirectURI,
		Endpoint:     provider.Endpoint(),
		Scopes:       cfg.Scopes,
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	return &UserService{
		repo:     repo,
		config:   oauth2Config,
		OIDC:     verifier,
		provider: provider,
	}
}
func (u *UserService) ExchangeCode(code string) (model.User, error) {
	token, err := u.config.Exchange(context.Background(), code)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to exchange code: %w", err)
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return model.User{}, fmt.Errorf("no id_token field in oauth2 token")
	}
	accessToken := token.AccessToken

	user, err := u.GetUserInfo(idToken, accessToken)
	if err != nil {
		return model.User{}, err
	}
	return user, nil

}

func (u *UserService) GetAuthURL(state string) string {
	return u.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (u *UserService) GetUserInfo(idToken, accessToken string) (model.User, error) {
	verifier := u.provider.Verifier(&oidc.Config{ClientID: u.config.ClientID})
	_, err := verifier.Verify(context.Background(), idToken)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to verify id token: %w", err)
	}
	userInfo, err := u.provider.UserInfo(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user info: %w", err)
	}
	var claims struct {
		PreferredUsername string `json:"preferred_username"`
		RealmAccess       struct {
			Roles []string `json:"roles"`
		} `json:"realm_access"`
		ResourceAccess map[string]struct {
			Roles []string `json:"roles"`
		} `json:"resource_access"`
	}
	if err := userInfo.Claims(&claims); err != nil {
		return model.User{}, fmt.Errorf("failed to parse id token claims: %w", err)
	}

	username := claims.PreferredUsername
	if username == "" {
		return model.User{}, fmt.Errorf("preferred_username not found in id token")
	}

	var roles []string
	roles = append(roles, claims.RealmAccess.Roles...)
	if clientRoles, exists := claims.ResourceAccess[u.config.ClientID]; exists {
		roles = append(roles, clientRoles.Roles...)
	}

	return model.User{
		Token:    idToken,
		Username: username,
		Roles:    roles,
	}, nil
}
