FROM golang:1.18-alpine

WORKDIR /app
COPY . .
WORKDIR /app/cmd/app

RUN go mod download

RUN go build -o /social-network

EXPOSE 3004

CMD [ "/social-network" ]
