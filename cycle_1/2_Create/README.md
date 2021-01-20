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

export FIREBASE_PROJECT="<my_firebase_project>"
export FIREBASE_DB_INSTANCE="<my_firebase_db_instance>"
```
---

Create gcloud configuration for Firebase Project
```bash
gcloud config configurations list

gcloud config configurations create ${FIREBASE_PROJECT}

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

Initialize Firebase CLI
```bash
firebase login
firebase init                         # choose the appropriate or the default resp.

firebase projects:list
firebase use ${FIREBASE_PROJECT}
```
---

Create a database instance
```bash
firebase database:instances:create ${FIREBASE_DB_INSTANCE}
firebase database:instances:list

firebase projects:addfirebase
```

Check respectively change Firestore rules
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
  
export GOOGLE_APPLICATION_CREDENTIALS="${LOCAL_CREDENTIALS_DIR}/${GCP_PROJECT}-${GCP_SA_NAME}.json"
ls -l $GOOGLE_APPLICATION_CREDENTIALS
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

--------
------


Change Firestore rules
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

old rules
```text
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /{document=**} {
      allow read, write: if request.auth.uid == "114690088071785841349";
    }
  }
}
```








