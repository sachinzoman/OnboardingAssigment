FROM docker-compose_php-apache

COPY . /var/www/html/
WORKDIR /var/www/html/
RUN mkdir -p restaurantgen  && chmod +x restaurant_proto_gen.sh && ./restaurant_proto_gen.sh
RUN ls && echo "#######################"
RUN cd restaurantgen && ls
RUN composer install
EXPOSE 80