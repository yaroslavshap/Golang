package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/http2"
)

func main() {
	folder := "/Users/aroslavsapoval/myProjects/data/images"
	files, err := filepath.Glob(folder + "/*")
	if err != nil {
		fmt.Println("Ошибка при чтении директории:", err)
		return
	}

	// Загрузите свой самоподписанный сертификат
	cert, err := os.ReadFile("/Users/aroslavsapoval/Desktop/Golang/My_server/server.cert")
	if err != nil {
		log.Fatalf("Couldn't load self-signed certificate: %s", err)
	}

	// Добавьте его в пул доверенных сертификатов
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("Couldn't append self-signed certificate to CertPool")
	}

	// Создайте новый http.Client с конфигурацией TLS, которая доверяет вашему самоподписанному сертификату
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: true,
		},
	}

	// Явно включите HTTP/2.
	http2.ConfigureTransport(tr)

	client := &http.Client{Transport: tr}

	// Остальной код...

	totalTime := 0.0
	for _, file := range files {
		fileData, err := os.Open(file)
		if err != nil {
			fmt.Println("Ошибка при открытии файла:", err)
			return
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(file))
		if err != nil {
			fmt.Println("Ошибка при создании формы файла:", err)
			return
		}
		io.Copy(part, fileData)

		writer.Close()
		req, err := http.NewRequest("POST", "https://127.0.0.1:8005/receive_images/", body)
		if err != nil {
			fmt.Println("Ошибка при создании запроса:", err)
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		startTime := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		resp.Body.Close()
		endTime := time.Now()

		transferTime := endTime.Sub(startTime).Seconds()
		totalTime += transferTime

		fmt.Printf("Response from server: %s, HTTP version: %s, %f секунд\n", resp.Status, resp.Proto, transferTime)
	}

	fmt.Printf("Общее время - %f\n", totalTime)
}
