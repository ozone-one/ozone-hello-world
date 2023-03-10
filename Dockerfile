FROM registry.ozone.one/library/nginx:mainline-alpine
RUN rm /etc/nginx/conf.d/*
ADD hello.conf /etc/nginx/conf.d/
ADD index.html /usr/share/nginx/html/
ADD logo512.png /usr/share/nginx/html/
ADD favicon.ico /usr/share/nginx/html/
