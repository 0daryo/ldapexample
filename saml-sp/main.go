package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"
)

func hello(w http.ResponseWriter, r *http.Request) {
	s := samlsp.SessionFromContext(r.Context())
	if s == nil {
		return
	}
	sa, ok := s.(samlsp.SessionWithAttributes)
	if !ok {
		return
	}
	fmt.Fprintf(w, "UserInfo\n")
	as := sa.GetAttributes()
	for k, v := range as {
		fmt.Printf("key: %s value: %s\n", k, v)
		fmt.Fprintf(w, "key: %s value: %s\n", k, v)
	}
}

func main() {
	confIDPMetadataURL := "http://localhost:18080/realms/master/protocol/saml/descriptor"
	rootURL, err := url.Parse("http://localhost:3000")
	if err != nil {
		log.Fatalf("failed to parse root URL: %v", err)
	}

	keyPair, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err)
	}
	idpMetadataURL, err := url.Parse(confIDPMetadataURL)
	if err != nil {
		panic(err)
	}
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient, *idpMetadataURL)
	if err != nil {
		panic(err)
	}
	samlSP, err := samlsp.New(samlsp.Options{
		EntityID:    "http://localhost:3000/saml/metadata",
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
		SignRequest: true,
	})
	if err != nil {
		log.Fatalf("failed to create samlsp: %v", err)
	}
	app := http.HandlerFunc(hello)
	http.Handle("/hello", samlSP.RequireAccount(app))
	http.Handle("/saml/", samlSP)
	http.ListenAndServe(":3000", nil)
}
