#!/bin/bash

# Build and run Odysseus terminal browser
echo "Building Odysseus..."
go build -o odysseus .

if [ $? -eq 0 ]; then
    echo "Starting Odysseus Terminal Browser..."
    echo "----------------------------------------"
    ./odysseus
else
    echo "Build failed!"
    exit 1
fi