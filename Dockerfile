FROM golang:1.18.2
WORKDIR /app
RUN go mod init myapi
COPY  *.go ./
RUN go mod tidy
RUN go mod download
RUN go build -o /api
CMD [ "/api" ]

#docker build -t myapi:1 .