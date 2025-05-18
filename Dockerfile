FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o wbrestapi

EXPOSE 3000

CMD [ "/app/wbrestapi" ]
