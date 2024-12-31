# Step 1: Build the app
FROM node:18 AS build

# Set the working directory
WORKDIR /app

# Install pnpm globally
RUN npm install -g pnpm

# Copy the package.json and pnpm-lock.yaml files
COPY package.json pnpm-lock.yaml ./

# Install dependencies using pnpm
RUN pnpm install

# Copy the rest of the application files
COPY . .

# Build the application using Vite
RUN pnpm run build

# Step 2: Serve the app with serve (node-based static server)
FROM node:18 AS production

# Install the "serve" package globally to serve the built app
RUN npm install -g serve

# Set the working directory
WORKDIR /app

# Copy the build output from the previous step
COPY --from=build /app/dist /app/dist

# Expose port 3000
EXPOSE 3000

# Start the app with serve (serve will serve the build folder on port 3000)
CMD ["serve", "-s", "dist", "-l", "3000"]