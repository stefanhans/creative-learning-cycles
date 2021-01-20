package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
)

//type ClcData struct {
//	Data       string `firestore:"data"`
//}

func main() {

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

	collection := client.Collection("clcdata")
	docRef, writeRes, err := collection.Add(ctx, map[string]interface{}{
		"data": "value",
	})
	if err != nil {
		log.Fatalf("Failed adding data: %v", err)
	} else {

		fmt.Printf("Success: document reference: %v\n", docRef)
		fmt.Printf("Success: write result: %v\n", writeRes)
	}

	//doc := collection.Doc("data")
	//
	//if err != nil {
	//	log.Fatalf("Failed to get document reference: %v", err)
	//}
	//
	//wr, err := doc.Create(ctx,  ClcData{
	//	Data: "create value",
	//})
	//
	//if err != nil {
	//	fmt.Printf("Failed create data: %v\n", err)
	//} else {
	//	fmt.Printf("Succeed create data: %v\n", wr)
	//}
	//
	//wr, err = doc.Set(ctx,  ClcData{
	//	Data: "set value",
	//})
	//
	//if err != nil {
	//	fmt.Printf("Failed set data: %v\n", err)
	//} else {
	//	fmt.Printf("Succeed set data: %v\n", wr)
	//}
}
