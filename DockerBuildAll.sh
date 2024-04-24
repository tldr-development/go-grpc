#!/bin/bash

TAG_BASE=${TAG_BASE}
VERSION=1.0.35
TARGETS=(account apns inspire)
PROJECT_ID=${PROJECT_ID}
APNS_URL=${$APNS_URL}

for TARGET in "${TARGETS[@]}"; do
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
        # dev
        if [ $TARGET == "inspire" ]; then
                for ENV in dev prod; do
                        gcloud run deploy $ENV-$TARGET-inspire \
                                --image=$TAG_BASE$TARGET:$VERSION \
                                --set-env-vars=ENV=$ENV,APP=inspire,PROJECT_ID=$PROJECT_ID,APNS_SERVER=$ENV$APNS_URL \
                                --region=asia-northeast3 \
                                --project=$PROJECT_ID \
                                && gcloud run services update-traffic $ENV-$TARGET-inspire --to-latest --region=asia-northeast3
                done
        fi

        if [ $TARGET == "account" ]; then
                for ENV in dev prod; do
                        gcloud run deploy $ENV-$TARGET-inspire \
                                --image=$TAG_BASE$TARGET:$VERSION \
                                --set-env-vars=ENV=$ENV,APP=inspire,PROJECT_ID=$PROJECT_ID \
                                --region=asia-northeast3 \
                                --project=$PROJECT_ID \
                                && gcloud run services update-traffic $ENV-$TARGET-inspire --to-latest --region=asia-northeast3
                done
        fi

        if [ $TARGET == "apns" ]; then
                for ENV in dev prod; do
                        gcloud run deploy $ENV-$TARGET-inspire \
                                --image=$TAG_BASE$TARGET:$VERSION \
                                --set-env-vars=ENV=$ENV,APP=inspire,PROJECT_ID=$PROJECT_ID \
                                --region=asia-northeast3 \
                                --project=$PROJECT_ID \
                                && gcloud run services update-traffic $ENV-$TARGET-inspire --to-latest --region=asia-northeast3
                done
        fi
done