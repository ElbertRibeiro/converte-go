package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func convertRTFtoPDF(rtfPath string, pdfPath string) error {
	// Verifique se o arquivo RTF existe
	if _, err := os.Stat(rtfPath); os.IsNotExist(err) {
		return fmt.Errorf("arquivo RTF não encontrado: %v", err)
	}

	// Comando para converter RTF para PDF usando LibreOffice
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf", rtfPath, "--outdir", filepath.Dir(pdfPath))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao converter RTF para PDF: %v\n%s", err, string(output))
	}

	// Verifique se o arquivo PDF foi criado
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return fmt.Errorf("arquivo PDF não foi criado: %v", err)
	}

	log.Printf("Arquivo PDF gerado com sucesso em: %s", pdfPath)
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
		http.Error(w, "Erro ao ler o arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Criar o diretório temporário se não existir
	tmpDir := "/tmp/rtf-to-pdf"
	err = os.MkdirAll(tmpDir, 0755)
	if err != nil {
		http.Error(w, "Erro ao criar diretório temporário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Criar um arquivo temporário para o RTF
	rtfPath := filepath.Join(tmpDir, header.Filename)
	pdfPath := filepath.Join(tmpDir, "output.pdf")

	outFile, err := os.Create(rtfPath)
	if err != nil {
		http.Error(w, "Erro ao criar arquivo temporário: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copiar o conteúdo do arquivo enviado para o arquivo temporário
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo temporário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Arquivo RTF salvo em: %s", rtfPath)
	log.Printf("Convertendo RTF para PDF...")

	// Converter RTF para PDF
	err = convertRTFtoPDF(rtfPath, pdfPath)
	if err != nil {
		log.Printf("Erro na conversão: %v", err)
		http.Error(w, "Erro ao converter RTF para PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Arquivo PDF gerado em: %s", pdfPath)

	// Abrir o arquivo PDF gerado
	pdfFile, err := os.Open(pdfPath)
	if err != nil {
		http.Error(w, "Erro ao abrir o arquivo PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer pdfFile.Close()

	// Definir o cabeçalho para indicar que o conteúdo é um PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=output.pdf")

	// Enviar o arquivo PDF como resposta
	_, err = io.Copy(w, pdfFile)
	if err != nil {
		http.Error(w, "Erro ao enviar o arquivo PDF: "+err.Error(), http.StatusInternalServerError)
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