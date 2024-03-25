package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/receive_images/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Protocols supported by the server: %v\n", r.Proto)

		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Fprintf(w, "Ошибка при получении файла: %v", err)
			return
		}
		defer file.Close()

		// err = os.MkdirAll("./new_images", os.ModePerm)
		// if err != nil {
		// 	fmt.Fprintf(w, "Ошибка при создании директории: %v", err)
		// 	return
		// }

		// out, err := os.Create("./new_images/" + header.Filename)
		// if err != nil {
		// 	fmt.Fprintf(w, "Ошибка при создании файла: %v", err)
		// 	return
		// }
		// defer out.Close()

		// _, err = io.Copy(out, file)
		// if err != nil {
		// 	fmt.Fprintf(w, "Ошибка при сохранении файла: %v", err)
		// 	return
		// }

		fmt.Fprintf(w, "Изображение %s успешно принято.", header.Filename)

		fmt.Printf("Изображение %s успешно принято.\n", header.Filename)
		fmt.Printf("Файл %s успешно сохранен.\n", header.Filename)
	})

	serverAddress := ":8005"
	fmt.Printf("Сервер запущен по адресу https://127.0.0.1%s\n", serverAddress)
	http.ListenAndServeTLS(serverAddress, "/Users/aroslavsapoval/Desktop/Golang/My_server/server.cert", "/Users/aroslavsapoval/Desktop/Golang/My_server/server.key", nil)
}
