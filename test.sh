#!/usr/bin/env bash
clear
cat data/shell/commands.txt | go run src/tsai.eu/orchestrator/orchestrator.go
