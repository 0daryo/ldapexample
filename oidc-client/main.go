package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	clientID     = "example-app"
	clientSecret = "example-secret"
	redirectURL  = "http://localhost:80/callback"
	providerURL  = "http://localhost:5556"

	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
)

func main() {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}
	verifier = provider.Verifier(&oidc.Config{ClientID: clientID})
	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// ハンドラ
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)

	fmt.Println("OIDC client running at: http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<a href="/login">Login with OIDC</a>`)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// stateとnonceは本来ランダムに生成
	state := "teststate"
	url := oauth2Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")
	oauth2Token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field", http.StatusInternalServerError)
		return
	}
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "ID Token verification failed", http.StatusInternalServerError)
		return
	}
	var claims struct {
		Email    string `json:"email"`
		Username string `json:"preferred_username"`
		Name     string `json:"name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
		return
	}
	fmt.Printf("claims: %+v\n", claims)
	fmt.Fprintf(w, "<h1>ログイン成功 🎉</h1><p>こんにちは、%s さん！</p>", claims.Name)
}
