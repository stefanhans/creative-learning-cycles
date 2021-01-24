package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
)

type ReadRequest struct {
	Path string `firestore:"path"`
	Type string `firestore:"type"`
	Data string `firestore:"data"`
}

type ReadResponse struct {
	Data string `firestore:"data"`
}

func main() {

	// Mock the request for reading a string data
	request := ReadRequest{
		Path: "creative-learning-cycles/my-doc-string",
		Type: "document",
	}

	// Mock the request for reading a non string data
	//request := ReadRequest{
	//	Path: "creative-learning-cycles/my-doc-no-string",
	//	Type: "document",
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

	//if request.Type == "collection" {
	//
	//	// Check if path contains an even number of IDs
	//	if strings.Count(request.Path, "/")%2 == 1 {
	//		fmt.Printf("Path %q has an even number of IDs", request.Path)
	//		return
	//	}
	//
	//	collection := client.Collection(request.Path)
	//	docRef, writeRes, err := collection.Add(ctx, map[string]interface{}{
	//		"data": request.Data,
	//	})
	//	if err != nil {
	//		log.Fatalf("Failed adding document: %v", err)
	//	} else {
	//
	//		fmt.Printf("Success: document reference: %v\n", docRef)
	//		fmt.Printf("Success: write result: %v\n", writeRes)
	//	}
	//}

	if request.Type == "document" {

		// Check if path contains an odd number of IDs
		if strings.Count(request.Path, "/")%2 == 0 {
			fmt.Printf("Path %q has an odd number of IDs", request.Path)
			return
		}

		doc := client.Doc(request.Path)
		snapshot, err := doc.Get(ctx)
		if err != nil {
			log.Fatalf("Failed doc.Get(ctx): %v", err)
		}

		// Convert data to expected struct
		var readResponse ReadResponse
		err = snapshot.DataTo(&readResponse)
		if err != nil {
			fmt.Println("failed to convert data to expected struct")

			data := snapshot.Data()
			fmt.Printf("data: %v (%T)\n", data, data)

			str, ok := data["data"].(string)
			if !ok {
				fmt.Println("failed assertion to string")
				fmt.Printf("value (type): %v (%T)\n", data["data"], data["data"])
			} else {
				fmt.Printf("value (type): %v (%T)\n", str, str)
			}
		} else {

			// Marshal struct to JSON
			bytes, err := json.MarshalIndent(readResponse, "", "  ")
			if err != nil {
				log.Fatalf("Failed json.MarshalIndent(readResponse,...): %v", err)
			} else {
				fmt.Printf("%v\n", string(bytes))
			}
		}
	}
}
