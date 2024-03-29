# Default to Go 1.19
ARG GO_VERSION=1.19

# Start from golang v1.15,6 base image
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group


# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates git

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import the code from the context.
COPY ./ ./

# Build the Go app (check this line)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix 'static' -o /app .
#RUN CGO_ENABLED=0 GOOS=arm64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main.


######## Start a new stage from scratch #######
# Final stage: the running container.
FROM scratch AS final
FROM --platform=linux/amd64 scratch
# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the compiled executable from the first stage.
COPY --chown=nobody:nobody --from=builder /app /app

# Perform any further action as an unprivileged user.
USER nobody:nobody

EXPOSE 7000
# Run the compiled binary.
ENTRYPOINT ["/app"]
