FROM nginx:stable-alpine as builder

# Set working directory for the build process
WORKDIR /app

# Copy the files for the frontend into the working directory
COPY ./frontend /app

# Install Node.js and npm
RUN apk add --no-cache nodejs npm

# Install frontend dependencies and perform build
RUN npm install
RUN npm run build

# Production stage
FROM nginx:stable-alpine

# Copy configuration file for nginx
COPY ./docker/nginx/frontend.conf /etc/nginx/conf.d/default.conf

# Copy the dist folder from the builder image to the nginx image
COPY --from=builder /app/dist /usr/share/nginx/html

# Open port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
