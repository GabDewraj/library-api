# - First Build
FROM golang:1.19.1-alpine as build_base


## Environment
WORKDIR /build
# Copy over all project files
COPY cmd /build/cmd
COPY pkgs /build/pkgs
COPY go.mod /build/go.mod
COPY go.sum /build/go.sum
RUN go mod tidy

# Build the app
RUN go build -o server "./cmd/main.go"

# Second build with Artifact from the previous build
FROM alpine:latest 
COPY --from=build_base /build/server .
COPY cmd/config/migrations cmd/config/migrations
#  Create a entry point to the executable.
ENTRYPOINT ["./server"]
