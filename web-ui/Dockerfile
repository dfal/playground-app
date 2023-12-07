FROM node:lts-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN cat version.txt
RUN sed -i "s/\"version\": \".*\"/\"version\": \"$(cat version.txt)\"/" package.json
RUN npm run build

FROM nginx:stable-alpine as release
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
