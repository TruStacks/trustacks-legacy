package server

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// toolchainsCollection is the name of toolchains firebase
// collection.
const toolchainsCollection = "toolchains"

// firebaseProvider is the firebase database provider.
type firebaseProvider struct {
	client *firestore.Client
}

// checkToolchainExists checks if the toolchain exists in the
// firebase collection.
func (p *firebaseProvider) checkToolchainExists(name string) (bool, error) {
	_, err := p.client.Collection(toolchainsCollection).Doc(name).Get(context.Background())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// createToolchain inserst the toolchain document in the firebase
// collection.
func (p *firebaseProvider) createToolchain(name string, data map[string]interface{}) error {
	_, err := p.client.Collection(toolchainsCollection).Doc(name).Create(context.Background(), data)
	return err
}

// newFirebaseProvider returns a firebase databse provider instance.
func newFirebaseProvider() (databaseProvider, error) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	fsClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error create the firestore client: %v", err)
	}
	return &firebaseProvider{fsClient}, nil
}
