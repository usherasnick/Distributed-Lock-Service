## Setup MySQL Cluster

1. create an internal Docker network that the containers will use to communicate

```shell
docker network create mysql-cluster --subnet=192.168.0.0/24
```

2. start one management node

```shell
docker run -d --net=mysql-cluster --name=management_n1 --ip=192.168.0.2 mysql/mysql-cluster:8.0 ndb_mgmd
```

3. start two data nodes

```shell
docker run -d --net=mysql-cluster --name=ndb_n1 --ip=192.168.0.3 mysql/mysql-cluster:8.0 ndbd
docker run -d --net=mysql-cluster --name=ndb_n2 --ip=192.168.0.4 mysql/mysql-cluster:8.0 ndbd
```

4. start the server node

```shell
docker run -d --net=mysql-cluster --name=mysql_n1 --ip=192.168.0.10 -e MYSQL_ROOT_PASSWORD='Pwd123!@' mysql/mysql-cluster:8.0 mysqld
```

5. start a container with an interactive management client to verify that the cluster is up

```shell
docker run -it --net=mysql-cluster mysql/mysql-cluster:8.0 ndb_mgm
```

6. enter the SHOW

```text
[Entrypoint] MySQL Docker Image 8.0.22-1.1.18-cluster
[Entrypoint] Starting ndb_mgm
-- NDB Cluster -- Management Client --
ndb_mgm> SHOW;
Connected to Management Server at: 192.168.0.2:1186
Cluster Configuration
---------------------
[ndbd(NDB)]	2 node(s)
id=2	@192.168.0.3  (mysql-8.0.22 ndb-8.0.22, Nodegroup: 0, *)
id=3	@192.168.0.4  (mysql-8.0.22 ndb-8.0.22, Nodegroup: 0)

[ndb_mgmd(MGM)]	1 node(s)
id=1	@192.168.0.2  (mysql-8.0.22 ndb-8.0.22)

[mysqld(API)]	1 node(s)
id=4	@192.168.0.10  (mysql-8.0.22 ndb-8.0.22)

ndb_mgm> 
```