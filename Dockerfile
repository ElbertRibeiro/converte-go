FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o rtf-to-pdf .

FROM elbertnilton/image-libreoffice

RUN mkdir -p /tmp/rtf-to-pdf && chmod 777 /tmp/rtf-to-pdf

COPY --from=builder /app/rtf-to-pdf /usr/local/bin/rtf-to-pdf

EXPOSE 8080

CMD ["rtf-to-pdf"]
