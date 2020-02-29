package service

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

var once sync.Once
var instance *firebase.App

type firebaseConfig struct {
	FirebaseType        string `json:"type"`
	ProjectID           string `json:"project_id"`
	PrivateKeyID        string `json:"private_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientID            string `json:"client_id"`
	AuthURL             string `json:"auth_url"`
	TokenURL            string `json:"token_url"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url"`
	ClientCertURL       string `json:"client_x509_cert_url"`
}

func initilizeFirebase() (*firebase.App, error) {
	config := firebaseConfig{
		FirebaseType:        os.Getenv("FIREBASE_CONFIG_TYPE"),
		ProjectID:           os.Getenv("FIREBASE_CONFIG_PROJECT_ID"),
		PrivateKeyID:        os.Getenv("FIREBASE_CONFIG_PRIVATE_KEY_ID"),
		PrivateKey:          os.Getenv("FIREBASE_CONFIG_PRIVATE_KEY"),
		ClientEmail:         os.Getenv("FIREBASE_CLIENT_EMAIL"),
		ClientID:            os.Getenv("FIREBASE_CLIENT_ID"),
		AuthURL:             os.Getenv("FIREBASE_AUTH_URL"),
		TokenURL:            os.Getenv("FIREBASE_TOKEN_URL"),
		AuthProviderCertURL: os.Getenv("FIREBASE_AUTH_PROVIDER_CERT_URL"),
		ClientCertURL:       os.Getenv("FIREBASE_CLIENT_CERT_URL"),
	}

	serializedConfig, _ := json.Marshal(config)
	opt := option.WithCredentialsJSON(serializedConfig)
	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		return nil, fmt.Errorf("Error initializing app: %v", err)
	}

	return app, nil
}

//GetFirebaseInstance get instance of firebase
func GetFirebaseInstance() (*firebase.App, error) {
	var err error

	if instance != nil {
		return instance, nil
	}

	once.Do(func() {
		instance, err = initilizeFirebase()
	})

	if err != nil {
		return nil, err
	}

	return instance, nil
}

//GetFirebaseAuth get firebase authentication client
func GetFirebaseAuth(context context.Context) (*auth.Client, error) {
	firebase, err := GetFirebaseInstance()

	if err != nil {
		return nil, fmt.Errorf("Cannot initialize firebase, %v", err)
	}

	auth, err := firebase.Auth(context)

	if err != nil {
		return nil, fmt.Errorf("Cannot initialize firebase authentication client, %v", err)
	}

	return auth, nil
}
