# Spice Paradise - Common Development Commands

This document lists common commands used for developing and running the Spice Paradise web application.

## I. Project Setup & Dependencies

1.  **Initialize Go Module (if `go.mod` doesn't exist):**
    *(Replace `spice.paradise/webapp` with your chosen module path)*
    ```bash
    go mod init spice.paradise/webapp
    ```

2.  **Download/Update Dependencies:**
    *(Ensures `go.mod` and `go.sum` are consistent and downloads packages)*
    ```bash
    go mod tidy
    ```

## II. Running the Application (Development Mode)

1.  **Set CGO_ENABLED Environment Variable (CRITICAL for `go-sqlite3`):**
    *This needs to be set in the terminal session where you run the Go commands.*

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

2.  **Run the Application:**
    *(This compiles and runs `main.go`)*
    ```bash
    go run main.go
    ```
    Access at `http://localhost:8080` (or as indicated in terminal logs).

## III. Building a Standalone Executable

1.  **Set CGO_ENABLED Environment Variable (as above).**

2.  **Build the Executable:**
    *(Creates `spice-paradise` or `spice-paradise.exe` in the project root)*
    ```bash
    go build -o spice-paradise .
    ```

3.  **Run the Built Executable:**
    *   **Linux/macOS:**
        ```bash
        ./spice-paradise
        ```
    *   **Windows:**
        ```bash
        .\spice-paradise.exe
        ```

## IV. Database Management (SQLite)

*The database file is `users.db` in the project root.*

1.  **Access via `sqlite3` Command Line Tool:**
    *(Requires `sqlite3` CLI to be installed and in PATH)*
    ```bash
    sqlite3 users.db
    ```

2.  **Common SQLite Commands (inside the `sqlite3` prompt):**
    *   List all tables:
        ```sql
        .tables
        ```
    *   Show schema for the `users` table:
        ```sql
        .schema users
        ```
    *   View all registered users:
        ```sql
        SELECT id, username, last_login_at FROM users;
        ```
    *   Count registered users:
        ```sql
        SELECT COUNT(*) FROM users;
        ```
    *   Exit `sqlite3`:
        ```sql
        .quit
        ```

3.  **Run a Single SQL Query from Terminal:**
    ```bash
    sqlite3 users.db "SELECT id, username FROM users;"
    ```

## V. Go Development Tools

1.  **Format Go Code:**
    *(Formats all `.go` files in the current directory and subdirectories)*
    ```bash
    go fmt ./...
    ```

2.  **Clean Go Build Cache:**
    *(Can help resolve stale build issues)*
    ```bash
    go clean -cache
    ```