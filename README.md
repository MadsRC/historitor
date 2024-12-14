# Historitor

[![Go Report Card](https://goreportcard.com/badge/github.com/MadsRC/historitor)](https://goreportcard.com/report/github.com/MadsRC/historitor)
[![Go Reference](https://pkg.go.dev/badge/github.com/MadsRC/historitor.svg)](https://pkg.go.dev/github.com/MadsRC/historitor)

![logo](./logo.png)  
*Logo created with Gopherkon at [quasilyte.dev](https://www.quasilyte.dev/gopherkon/?state=0e0k010v05090a0301020d030100000004)*

---

Historitor is a transactional log implementation, inspired by Redis and Kafka.

## Security

### Supply Chain

## Software Bill of Materials (SBOM)

An SBOM is generated for each release. The process is to create the release tag, push it, generate the SBOM and then create
a GitHub release for the version and attach the SBOM.

Creating the SBOM can be done with like this:

```shell
mise run generate-sbom
```

