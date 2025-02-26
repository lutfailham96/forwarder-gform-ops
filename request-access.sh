#!/bin/bash

echo "------------------------------------------------------------" >> request-access.log
./go-auto-commander >> request-access.log 2>&1
