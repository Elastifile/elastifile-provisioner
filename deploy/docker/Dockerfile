FROM fedora:24

USER root
RUN mkdir /prov
COPY elastifile-provisioner /prov/elastifile-provisioner

#RUN id

#COPY root.sh root.sh
#COPY rootshell /prov/rootshell
#RUN chmod 4777 /prov/rootshell
#RUN ./root.sh

#RUN echo "sleep 3600" > scr.sh
#RUN chmod 777 scr.sh
#COPY script.sh /tmp/script.sh

ENTRYPOINT ["/prov/elastifile-provisioner"]
#ENTRYPOINT "/tmp/script.sh"
