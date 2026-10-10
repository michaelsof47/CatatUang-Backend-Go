package config

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func InitFirebase() (*firebase.App, error) {
	// for testing in locally firebase emulators
	//conf := &firebase.Config{ProjectID: "catatuang"}

	conf := option.WithAuthCredentialsFile(option.ServiceAccount, "service_account.json")

	app, err := firebase.NewApp(context.Background(), nil, conf)
	if err != nil {
		return nil, err
	}

	return app, nil
}
