FROM python:3.10

RUN pip install pip --upgrade
RUN pip install ansible
RUN pip install psycopg2

RUN apt-get update -y && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    sshpass

RUN mkdir /.ansible
RUN chmod -R ug+rwx /.ansible

WORKDIR /ansible

RUN mkdir -p /ansible/roles
RUN mkdir -p /ansible/group_vars

COPY group_vars/all /ansible/group_vars/all
COPY group_vars/local/vars.yaml /ansible/group_vars/local
COPY roles/create_tables /ansible/roles/create_tables

COPY production /ansible
COPY pgo.yaml /ansible

CMD ["/bin/bash"]
