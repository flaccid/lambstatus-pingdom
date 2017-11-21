package main

import "github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"

func Handle(evt interface{}, ctx *runtime.Context) (string, error) {
	return "Hello, World!", nil
}
