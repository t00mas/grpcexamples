#!/bin/bash

# Set the default number of clients to 1 if no command-line argument is provided
num_clients=${1:-1}

# Run the go command for each client
for ((i=1; i<=$num_clients; i++)); do
  go run main.go -app=client &
done

# Wait for all client processes to finish
wait
