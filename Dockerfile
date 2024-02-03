FROM golang:latest AS build-stage
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./src

FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /app/app /app
COPY default_phrases.txt /default_phrases.txt
ENV HOST "0.0.0.0"
ENV PORT "80"
EXPOSE 80
USER nonroot:nonroot
ENTRYPOINT ["/app"]