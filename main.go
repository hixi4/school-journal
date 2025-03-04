package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

// Структура для учня
type Student struct {
	ID     string             `json:"id"`
	Name   string             `json:"name"`
	Grades map[string]float64 `json:"grades"`
}

// Структура для класу
type Class struct {
	Name     string    `json:"name"`
	Students []Student `json:"students"`
	Teacher  string    `json:"teacher"`
}

// Зберігання інформації про клас у пам'яті
var schoolClass = Class{
	Name: "10-A",
	Students: []Student{
		{ID: "1", Name: "Ivan", Grades: map[string]float64{"Math": 4.5, "History": 3.7}},
		{ID: "2", Name: "Maria", Grades: map[string]float64{"Math": 4.9, "History": 4.1}},
	},
	Teacher: "teacher123",
}

// Middleware для перевірки авторизації
func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != schoolClass.Teacher {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Обробник для отримання загальної інформації про клас
func getClassInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schoolClass)
}

// Обробник для отримання інформації про конкретного учня
func getStudentInfo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/student/")
	for _, student := range schoolClass.Students {
		if student.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(student)
			return
		}
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}

func main() {
	// Створення окремих шляхів для кращої маршрутизації
	http.Handle("/class", authorizationMiddleware(http.HandlerFunc(getClassInfo)))
	http.Handle("/student/", authorizationMiddleware(http.HandlerFunc(getStudentInfo)))

	fmt.Println("Сервер працює на динамічному порту")
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Слухає порт %d\n", listener.Addr().(*net.TCPAddr).Port)
	log.Fatal(http.Serve(listener, nil))
}
