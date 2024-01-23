FROM golang:1.21.4-alpine AS build_base
# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./s3-upload-server .

# Start fresh from a smaller image
FROM alpine
WORKDIR /app

ENV GIN_MODE=release
ENV PORT=8080

ENV SECRET_ACCESS_KEY=your-secret-access-key
ENV ACCESS_KEY_ID=your-access-key-id
ENV REGION=your-region
ENV BUCKET=your-bucket
ENV CDN_URL=your-cdn-url

#ENV API_ROUTER=/upload
#ENV PATH_PREFIX=/static



COPY --from=build_base /app/s3-upload-server /app/s3-upload-server
COPY --from=build_base /app/bin /app/bin

# Run the binary program produced by `go install`
CMD ["/app/s3-upload-server"]