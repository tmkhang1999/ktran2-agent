# Pull golang image from Dockerhub
FROM golang:latest

# Set up the working directory
WORKDIR /app

# copy the source code, then run build command
COPY . .
RUN go build .

# Execute the build
CMD ["./CSC482"]