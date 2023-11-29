FROM golang:1.21 as build
WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /dist

FROM build as test
RUN go test -v ./...

FROM alpine as release
WORKDIR /
COPY --from=build /dist /dist
RUN adduser -D nonroot
USER nonroot
ENV APP_VER=0.0.1
ENTRYPOINT [ "/dist" ]
