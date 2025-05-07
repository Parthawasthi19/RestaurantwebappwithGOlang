package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

// Structs for Cart Order Data

type BookingData struct {
	Name  string
	Email string
	Phone string
	Date  string
	Time  string
}

type CartItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type CartOrderData struct {
	Address  string     `json:"address"`
	Tip      float64    `json:"tip"`
	Subtotal float64    `json:"subtotal"`
	Total    float64    `json:"total"`
	Items    []CartItem `json:"items"`
}

// Handlers for different pages

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func serveMenuPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/menu.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func serveContactPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/contact.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func serveBookingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		booking := BookingData{
			Name:  r.FormValue("name"),
			Email: r.FormValue("email"),
			Phone: r.FormValue("phone"),
			Date:  r.FormValue("date"),
			Time:  r.FormValue("time"),
		}
		log.Printf("New Booking: %+v\n", booking)

		tmpl, err := template.ParseFiles("templates/booking.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, struct {
			BookingData
			ConfirmationMessage string
		}{
			BookingData:         booking,
			ConfirmationMessage: "Your table has been successfully booked! We look forward to your visit.",
		})
		return
	}

	tmpl, err := template.ParseFiles("templates/booking.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func serveCartPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/cart.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Handler to receive cart orders

func receiveCartOrder(w http.ResponseWriter, r *http.Request) {
	// Ensure we only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Unmarshal JSON into CartOrderData
	var order CartOrderData
	if err := json.Unmarshal(body, &order); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 🌟 Detailed logging of the order
	log.Println("🚚 New Order Received:")
	log.Printf("Address: %s\n", order.Address)
	log.Printf("Tip: ₹%.2f\n", order.Tip)
	log.Printf("Subtotal: ₹%.2f\n", order.Subtotal)
	log.Printf("Total: ₹%.2f\n", order.Total)
	log.Println("Items:")
	for _, item := range order.Items {
		log.Printf("- %s (Qty: %d) ₹%.2f each\n", item.Name, item.Quantity, item.Price)
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order received successfully!"))
}

// Serve static files (CSS, images, etc.)

func serveStaticFiles() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func main() {
	// Serve static files (CSS, images, etc.)
	serveStaticFiles()

	// Set up routes for different pages
	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/menu", serveMenuPage)
	http.HandleFunc("/contact", serveContactPage)
	http.HandleFunc("/booking", serveBookingPage)
	http.HandleFunc("/cart", serveCartPage)
	http.HandleFunc("/submit-cart", receiveCartOrder) // Endpoint for receiving orders

	// Start the server
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
