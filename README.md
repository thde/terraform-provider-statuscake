# StatusCake Terraform Provider

![tests](https://github.com/thde/terraform-provider-statuscake/workflows/test/badge.svg)
![golangci-lint](https://github.com/thde/terraform-provider-statuscake/workflows/golangci-lint/badge.svg)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13.x
- [Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

## Using the provider

The latest Docs can be found on the [Terraform Registry](https://registry.terraform.io/providers/thde/statuscake/latest/docs).

```hcl
terraform {
  required_providers {
    statuscake = {
      source = "thde/statuscake"
      version = "A.B.C"
    }
  }
}

provider "statuscake" {
  # Configuration options
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org).

### Build

```sh
$ cd terraform-provider-statuscake
$ make build
```

### Tests

```sh
$ cd terraform-provider-statuscake
$ make test
```

### Lint

```sh
$ cd terraform-provider-statuscake
$ make lint
```
