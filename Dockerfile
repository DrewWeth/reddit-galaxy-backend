FROM golang:alpine

WORKDIR /app
ADD subreddits.csv .


ADD main.go .
RUN go build

EXPOSE 8080

RUN ls
CMD ./app

