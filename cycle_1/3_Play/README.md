## Play

#### Limit the access to the Firestore database to the service account:

Get the client id from the credential file:
```bash
grep client_id ${LOCAL_CREDENTIALS_DIR}/${GCP_PROJECT}-${GCP_SA_NAME}.json
```

Change Firestore rules in the Firestore console:
```text
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /{document=**} {
      allow read, write: if request.auth.uid == "<client_id>";
    }
  }
}
```
---

#### Enhance the request to control the path and type, i.e. collection or document

```json
{
  "path": "creative-learning-cycles/cycle_1/",
  "type": "collection",
  "data": "some data"
}
```

Let's find out how using `main.go` copied from `../2_Create/main.go`.

---

Create new Go function in file `functions/store-data.go` using `template.txt` and `main.go`.

---

Prepare deployment and deploy Cloud Function `store-data`

```bash
cd functions

go mod init
go mod vendor
```

```bash
gcloud functions deploy store-data --region "${GCP_REGION}" \
    --entry-point StoreData --runtime go113 --trigger-http \
    --service-account="${GCP_SA_NAME}@${GCP_PROJECT}.iam.gserviceaccount.com" \
    --set-env-vars=GCP_PROJECT="${GCP_PROJECT}" \
    --update-labels=topic=creative-learning-cycles,cycle=1 \
    --allow-unauthenticated 
```

Get the URL and test the function

```bash
gcloud functions describe store-data --region "${GCP_REGION}" --format='value(httpsTrigger.url)'

curl <url> -d '{
    "path": "creative-learning-cycles/my-doc",
    "type": "document",
    "data": "setting data to a doc by function"
}'

curl <url> -d '{
    "path": "creative-learning-cycles/my-doc/my-new-collection",
    "type": "collection",
    "data": "adding a document to a new collection of a document by function"
}'

curl <url> -d '{
    "path": "creative-learning-cycles/cycle_1/my-collection",
    "type": "collection",
    "data": "adding a document to a collection by function"
}'

curl <url> -d '{
    "path": "creative-learning-cycles/cycle_1",
    "type": "document",
    "data": "setting data to a document with collection by function"
}'
```

#### Implement a read functionality


Let's find out how using `main.go` - save the old one under `main/path-and-type.go`

 
---

Next [Share](cycle_1/4_Share/README.md)