#!/bin/sh

PROC=50

for i in $(seq $PROC)
do
    ./injectu.py &
done
