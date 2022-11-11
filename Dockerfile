FROM golang:1.19-alpine3.16 as build

WORKDIR "/app"

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o rerun-workflow-for-branch

FROM alpine:3.16.2 as prod

RUN apk --no-cache add ca-certificates
RUN mkdir -p /app/home
WORKDIR /app/home
COPY --from=build /app/rerun-workflow-for-branch .

ENTRYPOINT ["./rerun-workflow-for-branch"]

