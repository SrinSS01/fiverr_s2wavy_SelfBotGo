FROM golang:1.22.4
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o selfbot .
EXPOSE 8090
CMD ["./selfbot"]
