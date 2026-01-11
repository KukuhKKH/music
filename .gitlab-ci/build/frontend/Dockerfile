FROM ${IMAGE_NAMESPACE}/base-music-frontend:latest

COPY frontend/.output/public /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
