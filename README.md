# Spice Paradise - Restaurant Web Application

Spice Paradise is a Go-based web application for a restaurant. It allows users to browse the menu, manage a shopping cart, book tables, and place food orders after registering and logging in.

## Features

*   User registration with password hashing (bcrypt).
*   User login and session management via cookies.
*   SQLite database for storing user credentials.
*   Static file serving for CSS and images.
*   Dynamic HTML page rendering using Go templates.
*   Client-side cart functionality using `localStorage`.
*   Authenticated order submission.
*   Table booking form.

## Project Structure

SPICE-PARADISE/
├── main.go # Core Go application logic and HTTP handlers
├── templates/ # HTML templates
│ ├── partials/ # Reusable template partials (e.g., navigation)
│ │ └── nav.html
│ ├── booking.html
│ ├── cart.html
│ ├── contact.html
│ ├── index.html
│ ├── login.html
│ ├── menu.html
│ └── register.html
├── static/ # Static assets (CSS, JavaScript, images)
│ ├── css/
│ │ └── styles.css
│ └── img/
│ └── (various images for logo, menu items, etc.)
├── go.mod # Go module file (defines module and dependencies)
├── go.sum # Go module checksums
└── users.db # SQLite database file (auto-created on first run)


## Prerequisites

1.  **Go:** Version 1.18 or newer installed. ([golang.org/dl/](https://golang.org/dl/))
2.  **C Compiler:** Required by the `go-sqlite3` driver.
    *   **Windows:** MinGW-w64 (recommended via [MSYS2](https://www.msys2.org/)) or TDM-GCC. Ensure the `bin` directory containing `gcc.exe` is added to your system's PATH environment variable.
    *   **Linux:** `gcc` (e.g., install `build-essential` on Debian/Ubuntu, or `Development Tools` group on Fedora/RHEL).
    *   **macOS:** Xcode Command Line Tools (install with `xcode-select --install`).
3.  **Git:** For version control (optional but recommended).

## Setup & Running

1.  **Clone the Repository (if you have one):**
    ```bash
    git clone <your-repository-url>
    cd SPICE-PARADISE
    ```
    If you don't have a repository, ensure you are in the `SPICE-PARADISE` project directory.

2.  **Initialize Go Modules & Install Dependencies:**
    If `go.mod` is not present, initialize it (replace `spice.paradise/webapp` with your desired module path):
    ```bash
    go mod init spice.paradise/webapp
    ```
    Then, fetch dependencies:
    ```bash
    go mod tidy
    ```
    This will download packages listed in `go.mod`, such as:
    *   `github.com/mattn/go-sqlite3`
    *   `golang.org/x/crypto/bcrypt`
    *   `github.com/google/uuid`

3.  **Set CGO_ENABLED (Crucial for `go-sqlite3`):**
    Before running, ensure CGo is enabled in your terminal session, especially if it's not default or on Windows:
    *   **PowerShell (Windows):**
        ```powershell
        $env:CGO_ENABLED="1"
        ```
    *   **Command Prompt (Windows):**
        ```cmd
        set CGO_ENABLED=1
        ```
    *   **Bash (Linux/macOS):**
        ```bash
        export CGO_ENABLED=1
        ```
        (Often, this is the default on Linux/macOS if a C compiler is found.)

4.  **Run the Application:**
    ```bash
    go run main.go
    ```
    The server will start (usually on `http://localhost:8080`). Check the terminal output for the exact address and any log messages.

5.  **Access in Browser:**
    Open your web browser and navigate to `http://localhost:8080`.

## Building an Executable

1.  Ensure CGo is enabled (see step 3 in "Setup & Running").
2.  Build the binary:
    ```bash
    go build -o spice-paradise .
    ```
    This creates an executable file named `spice-paradise` (or `spice-paradise.exe` on Windows).
3.  Run the executable:
    *   Linux/macOS: `./spice-paradise`
    *   Windows: `.\spice-paradise.exe`

## Database

*   The application uses SQLite, and the database file (`users.db`) is created in the project root.
*   To inspect the database, you can use the `sqlite3` command-line tool or a GUI tool like [DB Browser for SQLite](https://sqlitebrowser.org/).

    **Example using `sqlite3` CLI:**
    ```bash
    sqlite3 users.db
    ```
    Inside the SQLite prompt:
    ```sql
    SELECT id, username FROM users;
    .quit
    ```

## Development Notes

*   **Static Assets:** Images for menu items should be placed in the `static/img/` directory and referenced with paths like `/static/img/your-image.jpg` in HTML/CSS.
*   **Templates:** All HTML files are located in the `templates/` directory. Shared components like the navigation bar are in `templates/partials/`.