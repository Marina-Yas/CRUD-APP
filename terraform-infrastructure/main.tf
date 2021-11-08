terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  region  = "us-west-2"
  // if null, will read ENV AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
  
}

resource "aws_instance" "app_server" {
  ami           = "ami-830c94e3"
  instance_type = "t2.micro"
  user_data = <<EOF
	#!/bin/bash
  docker run -p 8080:8080 43dfbd37e028
  EOF
  
  tags = {
    Name = "ExampleAppServerInstance"
  }

}