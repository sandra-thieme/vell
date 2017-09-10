#!/bin/sh

vagrant docker-run centos7 -- fpm-cook --quiet package /vagrant/recipe.rb
