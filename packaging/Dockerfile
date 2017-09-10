FROM centos:centos7
MAINTAINER Sandra Thieme <thieme@synyx.de>
RUN yum clean all
RUN yum install -y rubygems ruby-devel make gcc tar rpm-build curl
RUN gem install fpm-cookery --no-ri --no-rdoc --version 0.33.0
RUN yum clean all

