package server

import (
	"context"
	"log"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
)

// deleteCollection deletes the collection.
func deleteCollection(ctx context.Context, client *firestore.Client, ref *firestore.CollectionRef, batchSize int) {
	for {
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			break
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestFirebaseCheckToolchainExists(t *testing.T) {
	client, err := firestore.NewClient(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	collection := client.Collection(toolchainsCollection)
	defer deleteCollection(context.Background(), client, collection, 10)
	exists, err := (&firebaseProvider{client}).checkToolchainExists("test")
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, exists)

	// create the document and test that it exists
	_, err = collection.Doc("test").Create(context.Background(), map[string]interface{}{"name": "test"})
	if err != nil {
		t.Fatal(err)
	}
	exists, err = (&firebaseProvider{client}).checkToolchainExists("test")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, exists)
}

func TestFirebaseCreateToolchain(t *testing.T) {
	client, err := firestore.NewClient(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	collection := client.Collection(toolchainsCollection)
	defer deleteCollection(context.Background(), client, collection, 10)
	params := map[string]interface{}{"name": "test", "sso": "authentik"}
	if err := (&firebaseProvider{client}).createToolchain("test", params); err != nil {
		t.Fatal(err)
	}
	doc, err := collection.Doc("test").Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test", doc.Data()["name"].(string), "got an unexpected toolchain name")
	assert.Equal(t, "authentik", doc.Data()["sso"].(string), "got an unexpected toolchain sso provider")
}
