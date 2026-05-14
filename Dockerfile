FROM node:lts-alpine

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm install --omit=dev

COPY server.js db.json ./

EXPOSE 9000

CMD ["node", "server.js"]
