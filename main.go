package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type UserData struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Email       string `json:"email"`
	UserTickets int    `json:"usertickets"`
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tickets", ticketHandler)
	http.HandleFunc("/buy", buyTicketHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

// Home Page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

// Ticket Form
func ticketHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "tickets.html", nil)
}

func validateUserInput(FirstName string, LastName string, Email string, UserTickets uint, RemainingTickets uint) (bool, bool, bool) {
	isValidName := len(FirstName) >= 2 && len(LastName) >= 2
	isValidEmail := strings.Contains(Email, "@")
	isValidTicketNumber := UserTickets >= 1 && UserTickets <= RemainingTickets

	return isValidName, isValidEmail, isValidTicketNumber
}

func buyTicketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	email := r.FormValue("email")
	userTickets, err := strconv.Atoi(r.FormValue("usertickets"))
	if err != nil {
		templates.ExecuteTemplate(w, "tickets.html", map[string]string{"error": "invalid ticket number"})
		return
	}

	// Store user data
	ticket := UserData{FirstName: firstName, LastName: lastName, Email: email, UserTickets: userTickets}
	saveTicket(ticket)

	// Redirect to confirmation Page
	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}

func saveTicket(ticket UserData) {
	file, err := os.OpenFile("data/tickets.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	data, err := json.Marshal(ticket)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	file.Write([]byte("\n"))
}
