## First stage - Build stage
# Pull golang image from Dockerhub
# FROM golang:alpine AS builder
FROM public.ecr.aws/docker/library/golang:alpine3.16 AS builder

# Set up the working directory
WORKDIR /app1

# copy the source code, then run build command
COPY ./modules .
RUN go build -o weather .

## Second stage - Run stage
# FROM alpine:latest
FROM public.ecr.aws/docker/library/golang:alpine3.16

# Set up the working directory
WORKDIR /app2

# Copy the executable binary file, env file and config file from the last stage to the new stage
COPY --from=builder /app1/weather .
COPY config.yaml .

# Execute the build
CMD ["./weather"]
