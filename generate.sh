#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
protoc blog/blogpb/blog.proto --go_out=plugins=grpc:.

# Node.js Client
protoc calculator/calculatorpb/calculator.proto --js_out=library=calculator/calculator_client_node,binary:calculator/calculator_client_node/build/gen