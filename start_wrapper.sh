#!/bin/bash
set -m

sleep 2 # postgres isnt up yet

/vodstream streamingester

fg %1