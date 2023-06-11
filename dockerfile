# Starting from a base Go image
FROM golang:1.19-alpine AS build

# Setting the working directory
WORKDIR /app

# Copying source code into the container
COPY . .

# Build the Go application (go binary)
RUN go build -o myapp

# using a minimal Alpine image
FROM alpine:latest

# Seting the working directory
WORKDIR /app

# copying config file containing postgres db connection details
COPY config/config.local.json config/config.local.json

RUN mkdir logs 

ENV CONFIGPATH=/app/config/config.local.json ENV=dev LOGDIR=/app/logs PORT=9000

# Copying the executable (binary) from the build container to the final image
COPY --from=build /app/myapp .

# Expose the port the application listens on
EXPOSE 9000

CMD ["./myapp"]
