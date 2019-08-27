package service

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

var instance *firebase.App

func initilizeFirebase() (*firebase.App, error) {
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CONFIG")))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, nil
}

//GetFirebaseInstance get instance of firebase
func GetFirebaseInstance() (*firebase.App, error) {
	if instance != nil {
		return instance, nil
	}

	instance, err := initilizeFirebase()

	if err != nil {
		return nil, err
	}

	return instance, nil
}
