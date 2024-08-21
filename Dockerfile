FROM golang:latest

COPY ./ ./

RUN go mod download
RUN go build -o main.exe ./cmd/reviews/main.go

CMD [ "./main.exe" ]