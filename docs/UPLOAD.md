# Logo inspection API

POST multipart/form-data to `/api/v1/logos/inspect` with exactly one file field named `logo`.

Currently **PNG only**. Validates signature, dimensions, full decoding, 5 MiB file limit, and 16 million pixel limit. Returns JSON metadata without persisting files.

```sh
curl -F 'logo=@example.png' http://localhost:8080/api/v1/logos/inspect
```

SVG and WebP require separate safe decoders and will follow in subsequent PRs. Asset generation is not yet implemented.
