FROM golang:1.18 AS gobuild

WORKDIR /go/src/github.com/expressoa2/randommatch/core
COPY . .
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM golang:1.18 AS gorun
WORKDIR /app
RUN update-ca-certificates
COPY --from=gobuild  /go/src/github.com/expressoa2/randommatch/core/main .
COPY --from=gobuild  /go/src/github.com/expressoa2/randommatch/core/api .
CMD ["./main"]