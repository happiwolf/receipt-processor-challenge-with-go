FROM golang:1.16
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /receiptprocessor
EXPOSE 8080
CMD [ "/receiptprocessor" ]
