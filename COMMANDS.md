# Project Setup Guide: RestaurantwebappwithGOlang on a New Windows System

This guide outlines the steps to set up and run the `RestaurantwebappwithGOlang` project on a fresh Windows system.

## Prerequisites

You will need to install the following software:

1.  **Git:** For version control and cloning the repository.
2.  **Go:** The programming language and toolchain.
3.  **A C Compiler (GCC via MSYS2):** Required because the project uses `go-sqlite3`, which depends on CGo.

---

## Step 1: Install Git

1.  **Download Git:**
    *   Go to [https://git-scm.com/download/win](https://git-scm.com/download/win) and download the latest Windows installer.
2.  **Install Git:**
    *   Run the installer.
    *   You can mostly accept the default settings.
    *   Ensure "Git Bash Here" and "Git GUI Here" are selected (usually default).
    *   For "Adjusting your PATH environment," the recommended option "Git from the command line and also from 3rd-party software" is generally best.
3.  **Verify Installation:**
    *   Open a **new** PowerShell window (Search for "PowerShell" in the Start Menu and open it).
    *   Type the following command and press Enter:
        ```powershell
        git --version
        ```
    *   You should see output like `git version 2.xx.x.windows.x`.

---

## Step 2: Install Go

1.  **Download Go:**
    *   Go to the official Go downloads page: [https://go.dev/dl/](https://go.dev/dl/)
    *   Download the Windows installer (MSI file).
2.  **Install Go:**
    *   Run the MSI installer and follow the on-screen prompts.
    *   The default installation location is usually `C:\Program Files\Go`. The installer should automatically add `C:\Program Files\Go\bin` to your system's PATH environment variable.
3.  **Verify Installation & PATH:**
    *   Open a **new** PowerShell window (it's important to use a new window after installations that modify the PATH).
    *   Type the following command and press Enter:
        ```powershell
        go version
        ```
    *   You should see output like `go version go1.2x.x windows/amd64`.

---

## Step 3: Install C Compiler (GCC via MSYS2)

The project requires a C compiler for the `go-sqlite3` dependency. We'll use GCC via MSYS2.

1.  **Download MSYS2:**
    *   Go to the MSYS2 official website: [https://www.msys2.org/](https://www.msys2.org/)
    *   Download the MSYS2 installer (e.g., `msys2-x86_64-xxxxxxxx.exe`).
2.  **Install MSYS2:**
    *   Run the installer. The recommended default installation path is `C:\msys64`.
3.  **Update MSYS2 and Install GCC:**
    *   After installation, search for "MSYS2 MSYS" in your Start Menu and open it. This will open an MSYS2 terminal (bash-like).
    *   **Update package database and core packages:**
        ```bash
        pacman -Syu
        ```
        If prompted to close the terminal (`[Y/n] y`), do so. Reopen "MSYS2 MSYS" and run the next command to complete updates:
        ```bash
        pacman -Su
        ```
    *   **Install the MinGW-w64 GCC toolchain (UCRT64 version is recommended):**
        ```bash
        pacman -S mingw-w64-ucrt-x86_64-gcc
        ```
        Press `Enter` to select the default choice and `Y` to proceed with the installation.
4.  **Add MinGW-w64 GCC to Windows PATH:**
    *   The GCC compiler will typically be located in `C:\msys64\ucrt64\bin` (if MSYS2 was installed to `C:\msys64`).
    *   Add this directory to your Windows System PATH:
        1.  In the Windows search bar, type "environment variables".
        2.  Click on "Edit the system environment variables".
        3.  In the System Properties window, click the "Environment Variables..." button.
        4.  In the "System variables" section, find the variable named `Path` and select it.
        5.  Click "Edit...".
        6.  Click "New" and add the path: `C:\msys64\ucrt64\bin`
        7.  Click "OK" on all open dialog windows to save the changes.
5.  **Verify GCC Installation:**
    *   Open a **NEW** PowerShell window.
    *   Type the following command and press Enter:
        ```powershell
        gcc --version
        ```
    *   You should see output from GCC, including its version (e.g., `gcc.exe (RevX, Built by MSYS2 project) X.X.X`). If you get a "command not found" error, double-check your PATH variable and ensure you're using a new PowerShell window.

---

## Step 4: Clone and Set Up The Project

1.  **Open a NEW PowerShell window.**
2.  **Choose a directory for your projects.** For example:
    ```powershell
    mkdir C:\MyProjects
    cd C:\MyProjects
    ```
    *(You can use any directory you prefer, e.g., `C:\Users\YourUserName\Desktop\GoWork`)*
3.  **Clone the repository:**
    ```powershell
    git clone https://github.com/Parthawasthi19/RestaurantwebappwithGOlang.git
    ```
4.  **Navigate into the cloned project directory:**
    ```powershell
    cd RestaurantwebappwithGOlang
    ```
5.  **Set CGO_ENABLED (Optional but recommended for first run):**
    Go usually enables CGo by default if a C compiler is found. To be explicit for the current session:
    ```powershell
    $env:CGO_ENABLED="1"
    ```
6.  **Download Go module dependencies:**
    The project uses Go Modules and already has a `go.mod` file. **DO NOT run `go mod init`.**
    ```powershell
    go mod tidy
    ```
    This command ensures all necessary dependencies are downloaded and the `go.sum` file is updated.

---

## Step 5: Build and Run Your Project

You can either run the application directly for development or build a standalone executable.

### Option A: Run Directly (for development)

This compiles and runs the application in one step.

1.  Ensure you are in the project's root directory (`RestaurantwebappwithGOlang`).
2.  If you didn't set `$env:CGO_ENABLED="1"` in the current session, you might want to do so.
3.  Execute:
    ```powershell
    go run main.go
    ```

### Option B: Build an Executable

This creates a standalone `.exe` file.

1.  Ensure you are in the project's root directory.
2.  If you didn't set `$env:CGO_ENABLED="1"` in the current session, you might want to do so.
3.  Build the application:
    ```powershell
    go build -o spice-paradise .
    ```
    *   `-o spice-paradise` names the output executable `spice-paradise.exe`.
    *   `.` indicates that the package in the current directory should be built.
4.  Run the compiled executable:
    ```powershell
    .\spice-paradise.exe
    ```
    (In PowerShell, you can often just type `.\spice-paradise`)

---

## Troubleshooting Tips

*   **"command not found" (`git`, `go`, `gcc`):**
    *   Verify that the respective `bin` directories (`C:\Program Files\Git\cmd`, `C:\Program Files\Go\bin`, `C:\msys64\ucrt64\bin`) are correctly listed in your System PATH.
    *   **Always open a new PowerShell window** after making changes to environment variables.
*   **CGo related errors / `gcc not found` during `go run` or `go build`:**
    *   Confirm GCC installation (Step 3.5).
    *   Confirm the MSYS2 GCC path is in your System PATH.
    *   Ensure you're in a new PowerShell window.
    *   Try explicitly setting `$env:CGO_ENABLED="1"` before the `go` command.
*   **Firewall Prompts:** If your application listens on a network port, the Windows Defender Firewall might ask for permission. Allow access if necessary.
*   **File/Directory Permissions:** If you encounter errors related to creating files or directories, ensure your user account has the necessary write permissions in the project location.