FROM centos:7
RUN yum -y install java-11-openjdk.x86_64 wget.x86_64 unzip.x86_64 git.x86_64 \
    && mkdir /cicd /java /js && cd /cicd && mkdir git \
    && wget https://downloads.gradle-dn.com/distributions/gradle-7.4.2-bin.zip \
    && unzip gradle-7.4.2-bin.zip
COPY ./cicd.sh /cicd
COPY ./rizzrack /cicd
ENTRYPOINT ["/cicd/rizzrack"]