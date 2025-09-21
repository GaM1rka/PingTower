# React + Vite + TypeScript (dev)
FROM node:20-alpine

WORKDIR /app

# Install deps first for better caching
COPY package*.json ./
RUN npm ci

# (Optional in dev; compose will mount the source)
COPY . .

EXPOSE 5173 24678

# Vite must listen on 0.0.0.0 inside containers
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "5173"]