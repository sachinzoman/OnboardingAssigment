FROM php:7.1.30-apache

RUN curl -sS https://getcomposer.org/installer | php \
  && chmod +x composer.phar && mv composer.phar /usr/local/bin/composer

#Install github
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git

#Prerequisite for grpc and protocol buffer
RUN apt-get -y install autoconf automake libtool curl make g++
RUN apt-get update && apt-get install -y zip unzip libz-dev

RUN pecl install grpc-1.19.0
RUN docker-php-ext-enable grpc

RUN docker-php-ext-install bcmath

# RUN cp "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"
ENV ENV docker

EXPOSE 80

COPY . /var/www/html/
WORKDIR /var/www/html/
RUN composer install
