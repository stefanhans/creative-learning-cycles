## Imagine

We want to send a JSON object 

```json
{
  "data": "value"
}
```

via HTTPS request 

```bash
curl https://<url>/store-data \
  -d '{
        "data": "value"
    }'

```
to a Cloud Function in Go
 
 ```bash
gcloud functions deploy store-data ...
```
 and store it in Firestore.
 
---

Next [Create](../2_Create/README.md#create)
