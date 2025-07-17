package firebaseauth

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	App        *firebase.App
	AuthClient *auth.Client
}

func NewFirebaseAuth(ctx context.Context, serviceAccountKeyPath string) (*FirebaseAuth, error) {
	if os.Getenv("RAILWAY_ENVIRONMENT_ID") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println(".env file not found, using system env vars instead")
		}
	}

	jsonCreds := os.Getenv("FIREBASE_CONFIG_JSON")
	if jsonCreds == "" {
		return nil, fmt.Errorf("FIREBASE_CONFIG_JSON env var not set")
	}

	creds := []byte(jsonCreds)
	opt := option.WithCredentialsJSON(creds)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting auth client: %w", err)
	}

	return &FirebaseAuth{
		App:        app,
		AuthClient: authClient,
	}, nil
}

func (f *FirebaseAuth) verifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return f.AuthClient.VerifyIDToken(ctx, idToken)
}

func (f *FirebaseAuth) GetUserInfoByIDToken(ctx context.Context, idToken string) (*auth.UserRecord, error) {
	token, err := f.verifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	userRecord, err := f.AuthClient.GetUser(ctx, token.UID)
	if err != nil {
		return nil, err
	}
	return userRecord, nil
}

func (f *FirebaseAuth) ListAllUsers(ctx context.Context) ([]*auth.ExportedUserRecord, error) {
	var users []*auth.ExportedUserRecord
	var pageToken string

	for {
		iter := f.AuthClient.Users(ctx, pageToken)
		for {
			user, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}
		pageToken = iter.PageInfo().Token
		if pageToken == "" {
			break
		}
	}

	return users, nil
}
