# Stacks
- fiber
- gRPC
- Github action
- GCP Cloud Build
- GCP Cloud Storage
- GCP Cloud Run
- GCP Datastore

# TEST


### Example for fiber as a client to gRPC server.

A sample program to showcase fiber as a client to a gRPC server.

#### Endpoints

| Method | URL           | Return value |
| ------ | ------------- | ------------ |
| GET    | "/add/:a/:b"  | a + b        |
| GET    | "/mult/:a/:b" | a \* b       |

#### Output

```bash
-> curl http://localhost:3000/add/33445/443234
{"result":"476679"}
-> curl http://localhost:3000/mult/33445/443234
{"result":"14823961130"}
```

-----
# [Cloud Run](https://cloud.google.com/run)
Deploy and run serverless container application

## Run locally
```bash
go run main.go
```

## Deploy
Make sure you have the permission

#### Manually
```bash
gcloud builds submit . \
    --substitutions SHORT_SHA=$(git rev-parse --short HEAD)
```

#### Automated
Trigger currently supports source code from:
 - [Cloud Source Repositories](https://cloud.google.com/source-repositories)
 - [Github] (https://github.com)

 Learn more on [official docs](https://cloud.google.com/cloud-build/docs/automating-builds/create-manage-triggers)