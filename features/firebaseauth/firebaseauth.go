package firebaseauth

import (
	"context"
	"encoding/base64"
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


func NewFirebaseAuth(ctx context.Context) (*FirebaseAuth, error) {
    // Load .env locally if not on Railway
    if os.Getenv("RAILWAY_ENVIRONMENT_ID") == "" {
        if err := godotenv.Load(); err != nil {
            log.Println(".env file not found, using system env vars instead")
        }
    }

    var creds []byte
    var err error

    b64Creds := os.Getenv("FIREBASE_CONFIG_B64")
    if b64Creds != "" {
        creds, err = base64.StdEncoding.DecodeString(b64Creds)
        if err != nil {
            return nil, fmt.Errorf("failed to decode FIREBASE_CONFIG_B64: %w", err)
        }
    } else {
        jsonCreds := os.Getenv("FIREBASE_CONFIG_JSON")
        if jsonCreds == "" {
            return nil, fmt.Errorf("no firebase credentials found in env vars")
        }
        creds = []byte(jsonCreds)
    }

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
