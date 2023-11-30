# gRPC Examples

A simple client-server example with keepalive connections.

## Quickstart

```bash
./gen_protos.sh
./run_server.sh
```

```bash
./run_clients.sh numClients
```

## Debug & Learn

To show ping frames and GOAWAYs due to idleness, run either or both of:

```bash
./debug_run_server.sh
```

```bash
./debug_run_clients.sh numClients
```