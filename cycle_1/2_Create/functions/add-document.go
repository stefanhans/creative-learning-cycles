package functions

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Data struct {
	Data string `firestore:"data"`
}

func AddDocument(w http.ResponseWriter, r *http.Request) {

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
	fmt.Printf("Request: %s\n", body)

	// Unmarshal request body
	bytes := []byte(string(body))
	var dataRequest Data
	err = json.Unmarshal(bytes, &dataRequest)
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

	// Close client when done.
	defer client.Close()

	collection := client.Collection("clcdata")
	docRef, writeRes, err := collection.Add(ctx, dataRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to add data: %s", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Success: document reference: %v\n", docRef)
	fmt.Printf("Success: write result: %v\n", writeRes)

	// Response
	_, err = fmt.Fprintf(w, "Success: %s\n", docRef.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}

	//jsonUsers, err := json.MarshalIndent(users, "    ", "    ")
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("failed to marshal data 'jsonUsers': %s", err), http.StatusInternalServerError)
	//	return
	//}
	//
	//// Response
	//_, err = fmt.Fprintf(w, "{ \n    %q: %s\n}\n", "Users", string(jsonUsers))
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	//}
}
