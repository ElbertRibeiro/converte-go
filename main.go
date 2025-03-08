package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/unidoc/unioffice/v2/document"
)

func convertRTFtoPDF(rtfPath string, pdfPath string) error {
	// Abrir o arquivo RTF
	doc, err := document.Open(rtfPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo RTF: %v", err)
	}

	// Salvar como PDF
	err = doc.SaveToFile(pdfPath)
	if err != nil {
		return fmt.Errorf("erro ao salvar o arquivo PDF: %v", err)
	}

	return nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Verificar se o método é POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Ler o arquivo enviado no corpo da requisição
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao ler o arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Criar um arquivo temporário para o RTF
	rtfPath := "/tmp/" + header.Filename
	pdfPath := "/tmp/output.pdf"

	outFile, err := os.Create(rtfPath)
	if err != nil {
		http.Error(w, "Erro ao criar arquivo temporário", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copiar o conteúdo do arquivo enviado para o arquivo temporário
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo temporário", http.StatusInternalServerError)
		return
	}

	// Converter RTF para PDF
	err = convertRTFtoPDF(rtfPath, pdfPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao converter RTF para PDF: %v", err), http.StatusInternalServerError)
		return
	}

	// Abrir o arquivo PDF gerado
	pdfFile, err := os.Open(pdfPath)
	if err != nil {
		http.Error(w, "Erro ao abrir o arquivo PDF", http.StatusInternalServerError)
		return
	}
	defer pdfFile.Close()

	// Definir o cabeçalho para indicar que o conteúdo é um PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=output.pdf")

	// Enviar o arquivo PDF como resposta
	_, err = io.Copy(w, pdfFile)
	if err != nil {
		http.Error(w, "Erro ao enviar o arquivo PDF", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Configurar o endpoint para upload
	http.HandleFunc("/upload", uploadHandler)

	// Iniciar o servidor
	fmt.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
