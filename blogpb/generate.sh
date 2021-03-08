#!/bin/bash
#Code generate script for blog RPC
protoc ./blog.proto --go_out=plugins=grpc:.