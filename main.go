package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"golang.org/x/crypto/bcrypt"
)

// --- Global Variables ---
var db *sql.DB
var tpl *template.Template

const sessionCookieName = "spice_paradise_session"
const sessionDuration = 24 * time.Hour // Session lasts for 24 hours

// --- Structs ---

// User struct for database
type User struct {
	ID           int
	Username     string
	PasswordHash string
}

// SessionData stores information about an active user session
type SessionData struct {
	UserID    int
	Username  string
	ExpiresAt time.Time
}

// In-memory session store (for simplicity; consider Redis/database for production)
var sessions = make(map[string]SessionData)

// PageData is passed to templates
type PageData struct {
	UserSession         *SessionData // Current logged-in user
	Error               string
	Message             string
	Data                interface{} // For any other specific data the page needs
	Year                int         // For footer copyright year
	ConfirmationMessage string      // For booking confirmation
	Booking             BookingData // For booking form prefill or display
}

// Structs for Cart Order Data (from your existing code)
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
	UserID   int        `json:"userId,omitempty"` // To associate order with user
	Address  string     `json:"address"`
	Tip      float64    `json:"tip"`
	Subtotal float64    `json:"subtotal"`
	Total    float64    `json:"total"`
	Items    []CartItem `json:"items"`
}

// --- Database Functions ---
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "username" TEXT NOT NULL UNIQUE,
        "password_hash" TEXT NOT NULL
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
	log.Println("Database initialized and users table checked/created.")
}

// --- Password Hashing ---
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// --- Session Management Functions ---
func createSession(w http.ResponseWriter, userID int, username string) string {
	sessionID := uuid.NewString()
	expiresAt := time.Now().Add(sessionDuration)

	sessions[sessionID] = SessionData{
		UserID:    userID,
		Username:  username,
		ExpiresAt: expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		// Secure: true, // Uncomment in production (requires HTTPS)
		// SameSite: http.SameSiteLaxMode,
	})
	log.Printf("Session created for user %s (ID: %d) with session ID %s", username, userID, sessionID)
	return sessionID
}

func getSessionUser(r *http.Request) (*SessionData, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, fmt.Errorf("no session cookie found: %w", err)
	}

	sessionID := cookie.Value
	sessionData, ok := sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("invalid session ID")
	}

	if sessionData.ExpiresAt.Before(time.Now()) {
		delete(sessions, sessionID) // Clean up expired session
		return nil, fmt.Errorf("session expired")
	}
	return &sessionData, nil
}

func deleteSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		sessionID := cookie.Value
		delete(sessions, sessionID)

		http.SetCookie(w, &http.Cookie{
			Name:    sessionCookieName,
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
			MaxAge:  -1,
		})
		log.Printf("Session deleted: %s", sessionID)
	}
}

// --- Middleware ---
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionUser, err := getSessionUser(r)
		if err != nil {
			log.Printf("Auth middleware: %s. Redirecting to login.", err)
			deleteSession(w, r) // Clear potentially invalid cookie
			// For API calls, might return 401, for browser, redirect
			// Check if it's an API call (e.g., from JavaScript fetch)
			if r.Header.Get("Accept") == "application/json" || r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required. Please login."})
				return
			}
			http.Redirect(w, r, "/login?error=Please+login+to+access+this+page.", http.StatusSeeOther)
			return
		}
		log.Printf("Auth middleware: User %s (ID: %d) authenticated. Proceeding.", sessionUser.Username, sessionUser.UserID)
		// Add user to context if needed for downstream handlers, for now, it's okay.
		next.ServeHTTP(w, r)
	}
}

// --- Helper to render templates ---
func renderTemplate(w http.ResponseWriter, tmplName string, data PageData) {
	data.Year = time.Now().Year() // Add current year for footer
	err := tpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		log.Printf("Error executing template %s: %v", tmplName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// --- HTTP Handlers (New and Modified) ---

func newPageData(r *http.Request) PageData {
	userSession, _ := getSessionUser(r) // Ignore error, if nil, user is not logged in
	return PageData{
		UserSession: userSession,
		Error:       r.URL.Query().Get("error"),
		Message:     r.URL.Query().Get("message"),
		Year:        time.Now().Year(),
	}
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // Handle your old "/index" route if necessary
		http.NotFound(w, r)
		return
	}
	data := newPageData(r)
	renderTemplate(w, "index.html", data)
}

func serveMenuPage(w http.ResponseWriter, r *http.Request) {
	data := newPageData(r)
	renderTemplate(w, "menu.html", data)
}

func serveContactPage(w http.ResponseWriter, r *http.Request) {
	data := newPageData(r)
	renderTemplate(w, "contact.html", data)
}

func serveBookingPage(w http.ResponseWriter, r *http.Request) {
	data := newPageData(r)
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			data.Error = "Failed to parse form."
			renderTemplate(w, "booking.html", data)
			return
		}
		booking := BookingData{
			Name:  r.FormValue("name"),
			Email: r.FormValue("email"),
			Phone: r.FormValue("phone"),
			Date:  r.FormValue("date"),
			Time:  r.FormValue("time"),
		}
		log.Printf("New Booking: %+v\n", booking)
		// Here you would typically save the booking to a database
		data.Booking = booking
		data.ConfirmationMessage = "Your table has been successfully booked! We look forward to your visit."
	}
	renderTemplate(w, "booking.html", data)
}

func serveCartPage(w http.ResponseWriter, r *http.Request) {
	data := newPageData(r)
	// The cart content is managed client-side by localStorage.
	// Login is only strictly required for *submitting* the order.
	renderTemplate(w, "cart.html", data)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := newPageData(r)
	if data.UserSession != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther) // Already logged in
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			data.Error = "Username and password are required."
			renderTemplate(w, "register.html", data)
			return
		}
		// Basic validation (add more as needed)
		if len(password) < 6 {
			data.Error = "Password must be at least 6 characters long."
			renderTemplate(w, "register.html", data)
			return
		}

		var existingUserID int
		err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&existingUserID)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error checking username existence: %v", err)
			data.Error = "Database error. Please try again."
			renderTemplate(w, "register.html", data)
			return
		}
		if existingUserID > 0 {
			data.Error = "Username already taken."
			renderTemplate(w, "register.html", data)
			return
		}

		hashedPassword, err := hashPassword(password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			data.Error = "Error processing registration."
			renderTemplate(w, "register.html", data)
			return
		}

		_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, hashedPassword)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			data.Error = "Failed to register user."
			renderTemplate(w, "register.html", data)
			return
		}

		log.Printf("User %s registered successfully", username)
		http.Redirect(w, r, "/login?message=Registration+successful.+Please+login.", http.StatusSeeOther)
		return
	}
	renderTemplate(w, "register.html", data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := newPageData(r)
	if data.UserSession != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther) // Already logged in
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			data.Error = "Username and password are required."
			renderTemplate(w, "login.html", data)
			return
		}

		var user User
		var storedHashedPassword string
		err := db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &storedHashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				data.Error = "Invalid username or password."
			} else {
				log.Printf("Error fetching user: %v", err)
				data.Error = "Database error."
			}
			renderTemplate(w, "login.html", data)
			return
		}

		if checkPasswordHash(password, storedHashedPassword) {
			createSession(w, user.ID, user.Username)
			log.Printf("User %s (ID: %d) logged in successfully", user.Username, user.ID)
			// Redirect to welcome page or home page
			// http.Redirect(w, r, "/welcome", http.StatusSeeOther)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			data.Error = "Invalid username or password."
			renderTemplate(w, "login.html", data)
			return
		}
	}
	renderTemplate(w, "login.html", data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	deleteSession(w, r)
	log.Println("User logged out.")
	http.Redirect(w, r, "/?message=Successfully+logged+out.", http.StatusSeeOther)
}

// This handler is now protected by authMiddleware
func receiveCartOrder(w http.ResponseWriter, r *http.Request) {
	// authMiddleware has already run and verified the user.
	// We can get user details if needed from the session,
	// but for now, we just know they are authenticated.
	sessionUser, err := getSessionUser(r)
	if err != nil {
		// Should not happen if authMiddleware is working, but good to check
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var order CartOrderData
	if err := json.Unmarshal(body, &order); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Associate order with the logged-in user
	order.UserID = sessionUser.UserID

	log.Println("🚚 New Order Received:")
	log.Printf("User ID: %d, Username: %s\n", order.UserID, sessionUser.Username)
	log.Printf("Address: %s\n", order.Address)
	log.Printf("Tip: ₹%.2f\n", order.Tip)
	log.Printf("Subtotal: ₹%.2f\n", order.Subtotal)
	log.Printf("Total: ₹%.2f\n", order.Total)
	log.Println("Items:")
	for _, item := range order.Items {
		log.Printf("- %s (Qty: %d) ₹%.2f each\n", item.Name, item.Quantity, item.Price)
	}
	// TODO: Save order to database here

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order received successfully!"})
}

// Serve static files (CSS, images, etc.)
func serveStaticFiles() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println("Serving static files from /static/ directory.")
}

func main() {
	// Initialize Database
	initDB()
	defer db.Close()

	// Parse all templates once at startup
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	tpl = template.Must(tpl.ParseGlob("templates/partials/*.html")) // Ensure partials are parsed
	log.Println("HTML templates parsed.")

	// Serve static files
	serveStaticFiles()

	// Set up routes
	http.HandleFunc("/", serveHomePage) // Main landing page
	// http.HandleFunc("/index", serveHomePage) // If you want to keep /index too
	http.HandleFunc("/menu", serveMenuPage)
	http.HandleFunc("/contact", serveContactPage)
	http.HandleFunc("/booking", serveBookingPage)
	http.HandleFunc("/cart", serveCartPage) // Viewing cart is public

	// Auth routes
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	// http.HandleFunc("/welcome", authMiddleware(serveWelcomePage)) // If you have a specific welcome page

	// Protected routes
	http.HandleFunc("/submit-cart", authMiddleware(receiveCartOrder)) // PROTECTED: User must be logged in

	// Start the server
	port := "8080"
	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
