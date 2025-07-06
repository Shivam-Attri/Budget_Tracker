# Full-Stack Budget Tracker Application

Welcome to the Full-Stack Budget Tracker application! This project is a complete, production-grade financial management tool featuring a secure Go backend API and a modern, responsive React frontend.

This repository is structured as a monorepo containing two main projects:
- `backend/`: A RESTful API built with Go
- `frontend/`: A Single-Page Application (SPA) built with React and Vite

## ✨ Core Features

This application is built with a focus on security, scalability, and a great user experience.

### Backend (Go)

- **Secure Authentication**: State-of-the-art JWT implementation with access and refresh tokens
- **Data Security**: Field-level AES-256 encryption for sensitive data and secure password hashing with bcrypt
- **Advanced API**: Features include pagination, dynamic filtering, sorting, and IP-based rate limiting
- **Production Ready**: Built with graceful shutdown, a health check endpoint, structured logging, and flexible configuration via Viper

### Frontend (React)

- **Modern UI/UX**: A clean, responsive interface built with React, Vite, and Tailwind CSS
- **Interactive Dashboard**: At-a-glance financial summaries with charts powered by Recharts
- **Full CRUD Functionality**: Seamlessly create, read, update, and delete transactions, categories, and budgets
- **Robust State Management**: Utilizes React Context for global state, ensuring a predictable and maintainable application

For a more detailed list of features, please see the individual README files in the backend and frontend directories.

## 🚀 Running the Entire Application (Docker)

The easiest and recommended way to get started is by using Docker and Docker Compose. This will build and run the entire application stack—database, backend, and frontend—with a single command.

### Prerequisites

- Docker
- Docker Compose

### Setup & Execution

1. **Clone the repository:**
   ```bash
   git clone <your-repository-url>
   cd <repository-directory>
   ```

2. **Configure Backend Secrets:**
   ```bash
   cd backend
   cp .env.example .env
   ```
   Open the newly created `backend/.env` file and replace the placeholder values with secure, randomly generated strings:
   - For `JWT_SECRET`: `openssl rand -base64 32`
   - For `ENCRYPTION_KEY`: `openssl rand -base64 32` (must be 32 bytes for AES-256)

3. **Start the Application Stack:**
   ```bash
   # Return to the root directory
   cd ..
   # Build and start all services
   docker-compose up --build
   ```

   This command will:
   - Build the container images for both the Go backend and the React frontend
   - Start the PostgreSQL database, backend, and frontend containers in the correct order
   - Connect the services on a shared Docker network

4. **Access the Application:**
   - Frontend UI: http://localhost:3000
   - Backend API: http://localhost:8000

5. **Stopping the Application**
   ```bash
   docker-compose down
   ```

## 📂 Project Structure

The project is organized as a monorepo to keep the backend and frontend codebases separate but managed together:

```
/budget-tracker-app
├── backend/
│   ├── ... (Go API source code)
│   └── README.md  # Backend-specific details
│
├── frontend/
│   ├── ... (React app source code)
│   └── README.md  # Frontend-specific details
│
├── docker-compose.yml
└── README.md      # You are here
```

For detailed instructions on developing, testing, or deploying a specific part of the application, please refer to the README file within its respective directory.