# Budget Tracker – React Frontend

This directory contains the complete source code for the Budget Tracker web application front end. It is a modern, responsive, and feature-rich Single-Page Application (SPA) built with React and Vite, designed to interact with the Go Budget Tracker API.

---

## ✨ Features

- **Secure Authentication:** Full login/registration flow with JWT handling
- **Automatic Token Refresh:** Seamlessly refreshes expired access tokens using refresh tokens for a persistent user session
- **Interactive Dashboard:** At-a-glance financial summary with charts to visualize income vs. expenses
- **Full CRUD for Transactions:** Create, read, update, and delete transactions through an intuitive modal form
- **Dynamic Transaction List:** A paginated table with advanced filtering (by type, date range, amount) and sorting capabilities
- **Category Management:** Users can create, edit, and delete their own custom spending categories
- **Responsive Design:** A clean, modern UI built with Tailwind CSS that works beautifully on all screen sizes, from mobile to desktop
- **Global State Management:** Uses React Context for managing authentication state and shared UI state (like modals) across the application

---

## 🛠️ Tech Stack

- **Framework:** [React](https://react.dev/) 18
- **Build Tool:** [Vite](https://vitejs.dev/)
- **Styling:** [Tailwind CSS](https://tailwindcss.com/)
- **Charting:** [Recharts](https://recharts.org/)
- **Icons:** [Lucide React](https://lucide.dev/)

---

## 🚀 Getting Started

These instructions will guide you through setting up and running the frontend application locally for development.

### Prerequisites

- [Node.js](https://nodejs.org/) (version 18 or later)
- [npm](https://www.npmjs.com/) or [yarn](https://yarnpkg.com/)
- A running instance of the [Go Budget Tracker API](<link-to-your-backend-repo-or-readme>)

### Setup & Installation

1. **Navigate to the `frontend` directory:**

   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   This command will download all the necessary packages defined in `package.json`.

   ```bash
   npm install
   ```

3. **Set up Environment Variables:**
   Create a new file named `.env.development.local` in the `frontend` directory. This file will hold the URL for your local backend API.

   ```env
   # frontend/.env.development.local
   VITE_API_URL=http://localhost:8000
   ```
   > **Note:** The `VITE_` prefix is required by Vite to expose the variable to your application code.

### Running the Development Server

Once the setup is complete, you can start the local development server.

1. **Run the `dev` script:**

   ```bash
   npm run dev
   ```

2. **Open the application:**
   The terminal will display a local URL, typically `http://localhost:5173`. Open this URL in your web browser to see the application running.

   The application will be in "hot-reload" mode, meaning any changes you make to the source code will be instantly reflected in the browser without a full page refresh.

---

## 📦 Building for Production

When you are ready to deploy the application, you need to create an optimized production build.

1. **Run the `build` script:**

   ```bash
   npm run build
   ```
   This command bundles all the application code, optimizes it for performance, and outputs the result into a `dist` folder.

2. **Deployment:**
   The contents of the `dist` folder are all you need to deploy. It contains static HTML, CSS, and JavaScript files. You can host this folder on any static hosting service, such as:
   - **Vercel** (Recommended)
   - **Netlify**
   - AWS S3
   - GitHub Pages

   For a production deployment, remember to set the `VITE_API_URL` environment variable in your hosting provider's settings to point to your deployed Go API's public URL.

---

## 📂 Project Structure

- **`public/`**: Contains static assets that are copied directly to the build output
- **`src/`**: Contains all the React source code
  - **`App.jsx`**: The main application component where all other components, pages, and context providers are assembled
  - **`main.jsx`**: The entry point of the application, which renders the `App` component into the DOM
  - **`index.css`**: The main stylesheet where Tailwind CSS directives are imported
- **`.env.development.local`**: Environment variables for local development
- **`Dockerfile`**: Instructions for building a production container image using Nginx
- **`nginx.conf`**: Nginx configuration to correctly serve the single-page application
- **`package.json`**: Lists project dependencies and scripts
- **`tailwind.config.js` & `postcss.config.js`**: Configuration files for Tailwind CSS
- **`vite.config.js`**: Configuration file for the Vite build tool

