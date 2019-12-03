# Modified from https://medium.com/@pierreprinetti/the-go-1-11-dockerfile-a3218319d191

# First stage: build the executable.
FROM alpine:3.10.3 AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN apk add --no-cache ca-certificates git

# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS requests.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create and own the binary folder
WORKDIR /app

ARG VERSION

# Copy the prebuilt binary
COPY bin/${VERSION}/linux_amd64/authproxy .

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/app/authproxy"]
