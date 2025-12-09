# BUILDER STAGE
# turning app into static binary so it doesnt rely on system libraries

# image will be used to composource code 
FROM golang:1.21-alpine AS builder

# set working directory

WORKDIR /app

# copy go.mod and go.sum, if not, docker caches go mod download

COPY go.mod go.sum ./
RUN go mod download

# copes rest of source code
COPY . .

# build go app (Compiles to static binary):
# CGO_ENABLED=0: atatically linked binary
# -o /go-service: setting output name of executable
#  -ldflags ="-s -w": removes symbol tables and debug info to make binary size smalelr
RUN CGO_ENABLED=0 go build -ldflags = "-s -w" -o /go-service main.go

# stage 2: final image
# smallest base image, no dependenscies, pretty much empty

FROM scratch

# copying nexessary certificates from builder stage for HTTPS client calls.

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# copy only compiled binary from builder stage for HTTP client calls
COPY --from=builder /go-service /go-service

# expose port for application lsitens
EXPOSE 8080

# command to run exec when container starts
ENTRYPOINT ["/go-service"]










