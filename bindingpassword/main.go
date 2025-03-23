package main

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

func main() {
	l, err := ldap.DialURL("ldap://localhost:389")
	if err != nil {
		log.Fatalf("接続失敗: %v", err)
	}
	defer l.Close()
	username := "ryotaro"
	password := "1111"
	userDN := fmt.Sprintf("cn=%s,dc=example,dc=com", username)

	err = l.Bind(userDN, password)
	if err != nil {
		fmt.Println("認証失敗:", err)
		return
	}
	fmt.Println("認証成功")
}
