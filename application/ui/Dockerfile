FROM nginx
CMD ["nginx", "-g", "daemon off;"]
# FROM node:16-slim AS ionic-build
# WORKDIR /app
# COPY package*.json /app/
# RUN npm install
# COPY ./ /app/
# RUN npx ionic build --prod

# FROM nginx:alpine
# RUN rm -rf /usr/share/nginx/html/*
# COPY --from=ionic-build /app/www/ /usr/share/nginx/html/
# COPY nginx/default.conf /etc/nginx/conf.d/
# STOPSIGNAL SIGTERM
# CMD ["nginx", "-g", "daemon off;"]

