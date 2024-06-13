FROM golang:1.22.4
WORKDIR /backend
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o selfbot .
EXPOSE 8000
CMD ["./selfbot", "serve", "--http=0.0.0.0:8000"]
