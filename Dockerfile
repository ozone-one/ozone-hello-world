# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM ozregistry.azurecr.io/ozone-public-registry/ozoneprod/golang:1.19.2-alpine3.16 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

ARG DOCKER_NETRC

# RUN apk update && apk upgrade && \
#     apk add --no-cache bash git openssh

# RUN echo "${DOCKER_NETRC}"

# RUN echo "${DOCKER_NETRC}" > ~/.netrc

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM ozregistry.azurecr.io/ozone-public-registry/ozoneprod/alpine:latest

# RUN apk --no-cache add ca-certificates

WORKDIR /bin

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /bin/master

EXPOSE	3000

ENTRYPOINT [ "/bin/master" ]
