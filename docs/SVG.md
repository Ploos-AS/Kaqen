# SVG inspection (restricted subset)

`POST /api/v1/logos/inspect` now accepts PNG and a **restricted, static SVG subset**. SVG files must contain an SVG root with a positive `viewBox` up to 4096 × 4096. Only basic vector shapes, paths and groups with explicitly allowed presentation attributes are accepted. Scripts, styles, embedded images, references, directives, XML namespaces other than SVG, and external resources are rejected.

This is a validation-only milestone: **no SVG rendering or sanitization output** is provided. Do not serve uploaded SVG bytes directly to browsers. Full SVG compatibility and WebP decoding are deferred to later changes.
