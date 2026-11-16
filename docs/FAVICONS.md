# Favicon ZIP endpoint

`POST /api/v1/favicons` with multipart field `logo` (PNG or static WebP).

```sh
curl -F 'logo=@logo.png' http://localhost:8080/api/v1/favicons -o kaqen-favicons.zip
```

Exports PNG sizes 16, 32, 180, 192, 512, a multi-resolution ICO (16, 32, 48), `site.webmanifest`, and `head.html`. PNGs retain transparency; aspect ratio is preserved and centered within a square canvas. The manifest contains placeholder site names; edit these before deploying.

**Not yet included:** SVG rasterization, background-color controls, maskable safe-zone icons, user interface, persistent storage, or full brand kits. This API does not accept SVG input. Do not publish untrusted SVG uploads.

Input limits: 5 MiB and 16 megapixels. ZIP generation is performed in memory with a 32 MiB output ceiling. Future work should stream/bound memory more tightly and add operational timeouts.
