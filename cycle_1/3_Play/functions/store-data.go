package functions

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type RequestData struct {
	Path string `firestore:"path"`
	Type string `firestore:"type"`
	Data string `firestore:"data"`
}

func StoreData(w http.ResponseWriter, r *http.Request) {

	// Gets Google Cloud Platform project ID
	gcpProject := os.Getenv("GCP_PROJECT")
	if gcpProject == "" {
		http.Error(w, fmt.Sprintf("internal error: GCP_PROJECT environment variable missing"), http.StatusInternalServerError)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read request: %s", err), http.StatusInternalServerError)
		return
	}

	// Unmarshal request body
	bytes := []byte(string(body))
	var requestData RequestData
	err = json.Unmarshal(bytes, &requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot unmarshall JSON input: %s", err), http.StatusInternalServerError)
		return
	}

	// Get a Firestore client.
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, gcpProject)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create client: %s", err), http.StatusInternalServerError)
		return
	}

	if requestData.Type == "collection" {

		// Check if path contains an even number of IDs
		if strings.Count(requestData.Path, "/")%2 == 1 {
			http.Error(w, fmt.Sprintf("path %q has an even number of IDs", requestData.Path), http.StatusInternalServerError)
			return
		}

		collection := client.Collection(requestData.Path)
		docRef, writeRes, err := collection.Add(ctx, map[string]interface{}{
			"data": requestData.Data,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed adding document: %s", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Success: document reference: %v\n", docRef)
		fmt.Printf("Success: write result: %v\n", writeRes)
	}

	if requestData.Type == "document" {

		// Check if path contains an odd number of IDs
		if strings.Count(requestData.Path, "/")%2 == 0 {
			http.Error(w, fmt.Sprintf("path %q has an odd number of IDs", requestData.Path), http.StatusInternalServerError)
			return
		}

		doc := client.Doc(requestData.Path)
		writeRes, err := doc.Set(ctx, map[string]interface{}{
			"data": requestData.Data,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed setting data: %s", err), http.StatusInternalServerError)
			return
		} else {
			fmt.Printf("Success: write result: %v\n", writeRes)
		}
	}

	// Response
	_, err = fmt.Fprintf(w, "implement the response")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}
}
