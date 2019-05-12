# Developers Guide

## Start required service (with Docker)

### NATS messaging

```
docker run -d --name nats-server -p 4222:4222 -p 6222:6222 -p 8222:8222 nats
```

### etcd key/value store

```
export NODE1="192.168.209.143"

docker volume create --name etcd-data
export DATA_DIR="etcd-data"

docker run -d \
  -p 2379:2379 \
  -p 2380:2380 \
  --net=host \
  --volume=${DATA_DIR}:/etcd-data \
  --name etcd quay.io/coreos/etcd:latest \
  /usr/local/bin/etcd \
  --data-dir=/etcd-data --name node1 \
  --initial-advertise-peer-urls http://${NODE1}:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://${NODE1}:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster node1=http://${NODE1}:2380
```
