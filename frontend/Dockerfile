# frontend/Dockerfile

# --- Stage 1: Build Stage ---
# Use an official Node.js image to build the React application.
FROM node:18-alpine AS build

# Set the working directory inside the container.
WORKDIR /app

# Copy package.json and package-lock.json to leverage Docker layer caching.
COPY package*.json ./

# Install project dependencies.
RUN npm install

# Copy the rest of the application source code.
COPY . .

# Build the application for production.
# This creates an optimized 'dist' folder with static assets.
RUN npm run build

# --- Stage 2: Serve Stage ---
# Use a lightweight Nginx image to serve the static files.
FROM nginx:1.25-alpine

# Copy the built assets from the 'build' stage to the Nginx public HTML directory.
COPY --from=build /app/dist /usr/share/nginx/html

# Copy the custom Nginx configuration file.
# This is crucial for single-page applications (SPAs) to handle routing correctly.
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80, which is the default Nginx port.
EXPOSE 80

# The default Nginx command will start the server.
CMD ["nginx", "-g", "daemon off;"]
