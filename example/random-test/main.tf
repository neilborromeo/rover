terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "3.1.0"
    }
  }
}

provider "random" {}

variable "max_length" {
  default = 5
}

resource "random_integer" "pet_length" {
  min = 1
  max = var.max_length
}

resource "random_pet" "dog" {
  length = random_integer.pet_length.result
}

locals {
  random_dog = random_pet.dog.id
}

resource "random_pet" "bird" {
  length = random_integer.pet_length.result
  prefix = local.random_dog
}

resource "random_pet" "dogs" {
  count = 3
  length = random_integer.pet_length.result
}

module "random_cat" {
  source = "./random-name"

  max_length = "3"
}

output "random_cat_name" {
  description = "random_cat_name"
  value = module.random_cat.random_name
}

resource "random_pet" "birds" {
  for_each = {
    "billy" = 1
    "bob" = 2
    "jill" = 3
  }
  
  prefix = each.key
  length = each.value
}