FROM golang:1.18-alpine

WORKDIR /app
COPY . .
WORKDIR /app/cmd/app

RUN go mod download

RUN go build -o /social-network

#EXPOSE 8080

CMD [ "/social-network" ]
