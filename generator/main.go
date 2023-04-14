package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/websocket"
)

// ProgramParameters - структура параметров программы
type ProgramParameters struct {
	Limit int `json:"limit"` // количество чисел, которые надо сгенерировать
	GoNum int `json:"goNum"` // количество горутин, которые пыхтят над генерацией чисел
}

// GetMainPage - Функция получения главной страницы
func GetMainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html") // index.html - этот файл сюда кидается
}

// ReadBufferSize и WriteBufferSize определяют размер буфера ввода-вывода в байтах
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// wsEndpoint - Функция, которая устанавливает что /ws это url для websocket соединения
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // Указываю метод валидации ориджина (сайта с которого пришел запрос)
	// Для любого ориджина возвращаю true, значит подойдет любой сайт

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Connected")

	reader(ws)
}

// reader - Получает сообщение с websocket и отправляет сгенерированные значения клиенту
func reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var params ProgramParameters
		json.Unmarshal(p, &params) // конвертирует json строку в объект ProgramParameters
		numbers := make(chan int)
		generateNumbers(params.Limit, params.GoNum, numbers)

		for item := range numbers {
			fmt.Println(item)
			time.Sleep(time.Millisecond * 100)
			conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(item)))
		}
	}
}

func main() {
	http.HandleFunc("/", GetMainPage)
	http.HandleFunc("/ws", wsEndpoint)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Printf("Starting server at port 8080\n")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}