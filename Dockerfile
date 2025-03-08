# Use a imagem oficial do Go como base
FROM golang:1.20-alpine AS builder

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o código-fonte para o diretório de trabalho
COPY . .

# Baixe as dependências do Go
RUN go mod download

# Compile o aplicativo Go
RUN go build -o rtf-to-pdf .

# Use uma imagem leve do Alpine para o estágio final
FROM alpine:latest

# Copie o binário compilado do estágio anterior
COPY --from=builder /app/rtf-to-pdf /usr/local/bin/rtf-to-pdf

# Expõe a porta 8080
EXPOSE 8080

# Defina o comando padrão para executar o aplicativo
CMD ["rtf-to-pdf"]