#####################################################################
###    STAGE 1: Build client and server                           ###
#####################################################################

FROM golang:1.20-alpine AS builder

# Preps
RUN apk update

# Set working directory for the build
WORKDIR /build

# Add source files
COPY . .

RUN go work sync

# Build server binary
WORKDIR /build/server

RUN go mod tidy
RUN go build -o /builder/bin/server

# Build client
WORKDIR /build/client
RUN go mod tidy
RUN go build -o /builder/bin/client

############################################
###    STAGE 2: Run server               ###
############################################

FROM alpine:3.12 as server

# Preps
RUN apk update

# Copy server binary
COPY --from=builder /builder/bin/server /app/server
COPY --from=builder /build/docker/demo-server.sh /bin/demo.sh
RUN chmod +x /bin/demo.sh

ARG SERVER_PORT
ARG SERVER_ADDRESS
ARG LOG_LEVEL
ARG STORAGE_PATH

# Set default values
ENV SERVER_PORT=8080
ENV SERVER_ADDRESS=localhost
ENV LOG_LEVEL=debug
ENV STORAGE_PATH=/storage

# Run server
CMD ["/app/server"]


############################################
###    STAGE 3: Run client               ###
############################################

FROM alpine:3.12 as client

# Preps
RUN apk update

# Copy server binary
COPY --from=builder /builder/bin/client /app/client
COPY --from=builder /build/docker/demo-client.sh /bin/demo.sh

RUN chmod +x /bin/demo.sh

ARG SERVER_PORT
ARG SERVER_ADDRESS
ARG LOG_LEVEL
ARG STORAGE_PATH

# Set default values
ENV SERVER_PORT=8080
ENV SERVER_ADDRESS=localhost
ENV LOG_LEVEL=debug
ENV STORAGE_PATH=/app/files


# Run server
CMD ["/app/client"]

