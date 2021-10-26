FROM centos:7

RUN  yum install -y epel-release && yum install -y ansible && yum install -y postgresql-devel && yum install -y install python-psycopg2

RUN mkdir /.ansible
RUN chmod -R ug+rwx /.ansible

WORKDIR /ansible

RUN mkdir -p /ansible/roles
RUN mkdir -p /ansible/group_vars

COPY group_vars/all /ansible/group_vars/all
COPY group_vars/local /ansible/group_vars/local
COPY roles /ansible/roles

COPY production /ansible
COPY pgo.yaml /ansible

# RUN ansible-playbook pgo.yaml -i production -l local

CMD ["/bin/bash"]