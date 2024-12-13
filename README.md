Historitor
===

[![Go Report Card](https://goreportcard.com/badge/github.com/MadsRC/historitor)](https://goreportcard.com/report/github.com/MadsRC/historitor)
[![Go Reference](https://pkg.go.dev/badge/github.com/MadsRC/historitor.svg)](https://pkg.go.dev/github.com/MadsRC/historitor)

Historitor is a transactional log implementation, inspired by Redis and Kafka.

## Security

### Supply Chain

## Software Bill of Materials (SBOM)

A SBOM is generated prior to each commit into the `main` branch and is saved as `sbom.json` in the root of the
repository.

While it is the authors belief that the ideal lifecycle-stage to generate an SBOM is as close to the build time as
possible, the current implementation is to generate the SBOM at the time of the commit. This is due to the fact that the
module does not build any artifacts on account of being a library.

If you would like to generate the SBOM yourself, you can do so by running the following command after checking out the
repository:

```shell
mise run generate-sbom
```