#!/bin/bash

# Prompt for input if environment variables are empty
if [ -z "$TAG_BASE" ]; then
  read -p "Enter TAG_BASE: " TAG_BASE
fi

if [ -z "$VERSION" ]; then
  read -p "Enter VERSION: " VERSION
fi

if [ -z "$TARGET" ]; then
  read -p "Enter TARGET: " TARGET
fi

if [ -z "$PROJECT_ID" ]; then
  read -p "Enter PROJECT_ID: " PROJECT_ID
fi

if [ -z "$APNS_URL" ]; then
  read -p "Enter APNS_URL: " APNS_URL
fi

if [ -z "$ENV" ]; then
  read -p "Enter ENV (dev/prod): " ENV
fi

TAG_BASE=${TAG_BASE}
VERSION=${VERSION}
TARGET=${TARGET}
PROJECT_ID=${PROJECT_ID}
APNS_URL=${$APNS_URL}
ENV=${ENV}

echo "Building $TAG_BASE$TARGET:$VERSION"
echo "Project ID: $PROJECT_ID"
echo "APNS URL: $APNS_URL"
echo "Environment: $ENV"

# confirm 
echo "Confirming build for $TAG_BASE$TARGET:$VERSION"
read -p "Is this correct? (y/n): " confirm
if [[ $confirm != "y" ]]; then
  echo "Build cancelled."
  exit 1
fi


cat << EOF > Dockerfile
FROM golang:1.21 as build-go
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app /app/$TARGET

FROM alpine:latest
RUN addgroup -S app && adduser -S app -G app
USER app
WORKDIR /home/app
COPY --from=build-go /bin/app ./
EXPOSE 3000
ENTRYPOINT ["./app"]

EOF

echo "Building $TAG_BASE$TARGET:$VERSION"
# docker build -t $TAG_BASE$TARGET:$VERSION -f Dockerfile .
# push to gcr
gcloud builds submit . --tag $TAG_BASE$TARGET:$VERSION
rm Dockerfile
gcloud run deploy $ENV-$TARGET \
        --image=$TAG_BASE$TARGET:$VERSION \
        --set-env-vars=ENV=$ENV,APP=$TARGET,PROJECT_ID=$PROJECT_ID \
        --region=asia-northeast3 \
        --project=$PROJECT_ID \
        && gcloud run services update-traffic $ENV-$TARGET --to-latest --region=asia-northeast3