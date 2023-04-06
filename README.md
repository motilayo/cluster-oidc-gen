# Cluster OIDC Gen

Cluster OIDC Gen is a Go-based application that allows you to generate OpenID Connect (OIDC) configurations for clusters. It generates the discovery document and JWKS (JSON Web Key Set), and then uploads them to a specified public endpoint (e.g., AWS S3, Azure Storage Account Container, Google Cloud Storage Bucket).

**Currently, it only supports uploading to an AWS S3 bucket's public endpoint.**

*Support for Azure Storage Account Container, Google Cloud Storage Bucket will be added soon*

## Features

- Generates OIDC configurations for clusters
- Generates the discovery document and JWKS JSON files
- Uploads the generated files to a public endpoint

## Installation

To install Cluster OIDC Gen, you can use `go get`:
```sh
go get -u github.com/motilayo/cluster-oidc-gen
```


## Usage

Cluster OIDC Gen takes the following input parameter:

- `--bucket-name`: The name of the bucket in the cloud storage provider where the discovery document and JWKS JSON files will be uploaded.

Here's an example of how to use Cluster OIDC Gen:

```sh
cluster-oidc-gen --bucket-name <your-bucket-name>
```

## Contributing

Contributions are welcome! If you would like to contribute to Cluster OIDC Gen, please fork the repository, create a branch, make your changes, and submit a pull request.

## License

This is free and unencumbered public domain software. For more information, see http://unlicense.org/ or the accompanying [LICENSE](LICENSE) file.

