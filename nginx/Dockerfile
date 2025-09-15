FROM nginx

STOPSIGNAL SIGQUIT

EXPOSE 80 443

RUN rm -f /etc/nginx/conf.d/default.conf

COPY app.conf /etc/nginx/conf.d/app.conf
COPY certs/ /etc/nginx/certs/

RUN nginx -t
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD wget -qO- http://127.0.0.1/ || exit 1

CMD ["nginx", "-g", "daemon off;"]