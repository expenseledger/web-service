package service

import (
	"fmt"
	"os"
	"sync"

	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

var once sync.Once
var instance *firebase.App

func initilizeFirebase() (*firebase.App, error) {
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CONFIG")))
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
