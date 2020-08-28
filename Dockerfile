FROM golang:1.14
MAINTAINER syedparvez@gmail.com
WORKDIR /app
COPY . /app
RUN go get -d -v ./...
RUN go build app.go
#comma seperated email ids
ENV INFRA_MAILS_IDS syedparvez72@gmail.com
ENV GO_ENV production
ENV SENDGRID_API_KEY ""
# namespace where whatsapp running, required for increasing the pvc size
ENV WHATSAPP_ACCOUNT testwhatapp-namespace
ENV IN_CLUSTER true
VOLUME /wamedia
CMD ["./app"]