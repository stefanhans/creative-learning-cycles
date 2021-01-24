package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
)

type ClcData struct {
	Path string `firestore:"path"`
	Type string `firestore:"type"`
	Data string `firestore:"data"`
}

func main() {

	// Mock the request for setting data to a doc
	request := ClcData{
		Path: "creative-learning-cycles/my-doc",
		Type: "document",
		Data: "setting data to a doc",
	}

	// Mock the request for adding a document to a new collection of a document
	//request := ClcData{
	//	Path: "creative-learning-cycles/my-doc/my-new-collection",
	//	Type: "collection",
	//	Data: "adding a document to a new collection of a document",
	//}

	// Mock the request for adding a document to a collection
	//request := ClcData{
	//	Path: "creative-learning-cycles/cycle_1/my-collection",
	//	Type: "collection",
	//	Data: "adding a document to a collection",
	//}

	// Mock the request for setting data to a document with collection
	//request := ClcData{
	//	Path: "creative-learning-cycles/cycle_1",
	//	Type: "document",
	//	Data: "setting data to a document with collection",
	//}

	// Normalize path
	request.Path = strings.Trim(request.Path, "/")

	// Gets Google Cloud Platform project ID.
	gcpProject := os.Getenv("GCP_PROJECT")
	if gcpProject == "" {
		fmt.Printf("GCP_PROJECT environment variable missing")
		return
	}

	// Gets local directory for credentials.
	localCredentialDir := os.Getenv("LOCAL_CREDENTIALS_DIR")
	if localCredentialDir == "" {
		fmt.Printf("LOCAL_CREDENTIALS_DIR environment variable missing")
		return
	}

	// Gets Google Cloud Platform service account name.
	gcpSaName := os.Getenv("GCP_SA_NAME")
	if gcpSaName == "" {
		fmt.Printf("GCP_SA_NAME environment variable missing")
		return
	}

	// Credential file
	credentialFile := fmt.Sprintf("%s/%s-%s.json", localCredentialDir, gcpProject, gcpSaName)

	// Get a Firestore client.
	ctx := context.Background()

	// Set the credential file as an option
	sa := option.WithCredentialsFile(credentialFile)

	// Get a new client
	client, err := firestore.NewClient(ctx, gcpProject, sa)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Close client when done.
	defer client.Close()

	if request.Type == "collection" {

		// Check if path contains an even number of IDs
		if strings.Count(request.Path, "/")%2 == 1 {
			fmt.Printf("Path %q has an even number of IDs", request.Path)
			return
		}

		collection := client.Collection(request.Path)
		docRef, writeRes, err := collection.Add(ctx, map[string]interface{}{
			"data": request.Data,
		})
		if err != nil {
			log.Fatalf("Failed adding document: %v", err)
		} else {

			fmt.Printf("Success: document reference: %v\n", docRef)
			fmt.Printf("Success: write result: %v\n", writeRes)
		}
	}

	if request.Type == "document" {

		// Check if path contains an odd number of IDs
		if strings.Count(request.Path, "/")%2 == 0 {
			fmt.Printf("Path %q has an odd number of IDs", request.Path)
			return
		}

		doc := client.Doc(request.Path)
		writeRes, err := doc.Set(ctx, map[string]interface{}{
			"data": request.Data,
		})
		if err != nil {
			log.Fatalf("Failed setting data: %v", err)
		} else {
			fmt.Printf("Success: write result: %v\n", writeRes)
		}
	}
}
