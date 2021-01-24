## Create

#### Requirements

You need 

- a project on [Google Cloud Platform](https://console.cloud.google.com/) 
- a project on [Cloud Firebase](https://console.firebase.google.com/) with 
    - a Cloud Firestore database 
    - [firebase-tools](https://firebase.google.com/docs/cli)
---

Set environment variables

```bash
export GCP_PROJECT="<my_gcp_project>"
export GCP_REGION="<my_gcp_region>"
export GCP_ZONE="<my_gcp_zone>" # optional
export GCP_ACCOUNT="<my_gcp_account>"
export GCP_SA_NAME="<my_service_account_name>"

export LOCAL_CREDENTIALS_DIR="<local credentials directory>"
```
---

Create gcloud configuration
```bash
gcloud config configurations list

gcloud config configurations create ${GCP_PROJECT}

gcloud config configurations list
```
---

Configure gcloud configuration
```bash
gcloud config list

gcloud config set core/project ${GCP_PROJECT}
gcloud config set core/account ${GCP_ACCOUNT}
gcloud config set compute/region ${GCP_REGION}
gcloud config set compute/zone ${GCP_ZONE} # optional

gcloud config list
gcloud config configurations list
```
---

Initialize Firebase CLI for Firestore
```bash
firebase login
firebase init firestore # choose the appropriate or the default resp.

firebase projects:list
firebase use ${GCP_PROJECT}
```

---

Create the database by `gcloud`
```bash
gcloud firestore databases create --region=${GCP_REGION}
```

---

Check respectively change Firestore rules in `firestore.rules`
```text
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /{document=**} {
      allow read, write: if request.auth != null;
    }
  }
}
```

Deploy if changed
```bash
firebase deploy --only firestore:rules
```


#### Service Account

Create a service-account to connect to the Firestore database
```bash
gcloud iam service-accounts create ${GCP_SA_NAME} \
    --description="Service account to access Firestore API" \
    --display-name="${GCP_SA_NAME}"
    
gcloud projects add-iam-policy-binding ${GCP_PROJECT} \
    --member "serviceAccount:${GCP_SA_NAME}@${GCP_PROJECT}.iam.gserviceaccount.com" \
    --role "roles/firebasedatabase.admin"    
      
gcloud iam service-accounts keys create ${LOCAL_CREDENTIALS_DIR}/${GCP_PROJECT}-${GCP_SA_NAME}.json \
  --iam-account ${GCP_SA_NAME}@${GCP_PROJECT}.iam.gserviceaccount.com
```

### Cloud Function

Use `main.go` to test functionality in main function

```bash
go mod init

go run main.go
```

---

Create new Go function in file `functions/add-document.go` using `template.txt` and `main.go`.

---

Prepare deployment and deploy Cloud Function `add-document`

```bash
cd functions

go mod init
go mod vendor
```

```bash
gcloud functions deploy add-document --region "${GCP_REGION}" \
    --entry-point AddDocument --runtime go113 --trigger-http \
    --service-account="${GCP_SA_NAME}@${GCP_PROJECT}.iam.gserviceaccount.com" \
    --set-env-vars=GCP_PROJECT="${GCP_PROJECT}" \
    --update-labels=topic=creative-learning-cycles,cycle=1 \
    --allow-unauthenticated 
```

---

Get the URL and test the function

```bash
gcloud functions describe add-document --region "${GCP_REGION}" --format='value(httpsTrigger.url)'

curl <url> -d '{"data": "my value"}'

```

---

Next [Play](../3_Play/README.md#play)





