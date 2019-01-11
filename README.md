# Just S3

This is a CLI tool that copies to and from S3. It's statically-linked
Golang and it has no other dependencies.

I wrote it because usually all I want to do in my AMI or container or
instance is to grab a single file from S3, and I'm tired of installing
the entire `awscli` package, and its Python dependencies, just to do
that.

## Usage

```
# Copy from a local file to S3
just-s3 cp /my/local/file s3://bucket/key

# or, copy from S3 to a local file
just-s3 cp s3://bucket/key /tmp/local
```

Credentials will be taken from the usual AWS sources -- environment
variables, AWS config files, instance profiles, etc.

You must set the `AWS_DEFAULT_REGION` or `AWS_REGION` environment
variable.

## Caveats

S3-to-S3 copies are implemented using a temporary file, so you cannot
copy an S3 file larger than your filesystem will contain.
