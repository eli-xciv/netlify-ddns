FROM golang:1.16-alpine

RUN go install github.com/eli-xciv/netlify-ddns@0.0.1-alpha

ENTRYPOINT ["netlify-ddns"]
