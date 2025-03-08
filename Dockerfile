# Use a imagem oficial do Go como base
FROM golang:1.20-alpine AS builder

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o código-fonte para o diretório de trabalho
COPY . .

# Baixe as dependências do Go e gere o go.sum
RUN go mod tidy

# Compile o aplicativo Go
RUN go build -o rtf-to-pdf .

# Use uma imagem base com LibreOffice instalado
FROM alpine:latest

# Instale o LibreOffice e dependências necessárias (fontes TrueType)
RUN apk add --no-cache libreoffice ttf-dejavu

# Crie o diretório temporário e defina permissões
RUN mkdir -p /tmp/rtf-to-pdf && chmod 777 /tmp/rtf-to-pdf

# Copie o binário compilado do estágio anterior
COPY --from=builder /app/rtf-to-pdf /usr/local/bin/rtf-to-pdf

# Expõe a porta 8080
EXPOSE 8080

# Defina o comando padrão para executar o aplicativo
CMD ["rtf-to-pdf"]