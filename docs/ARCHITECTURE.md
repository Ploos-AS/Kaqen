# Architecture

Go standard-library HTTP service initially. `internal/api` owns routing; image validation, transformation and export will be separate packages.

Planned flow: upload -> bounded validation and safe decoding -> user-controlled normalization -> icon rendering -> ZIP -> download -> cleanup.

Security requirements before upload support: body limits, pixel limits, no SVG scripts or external resource loading, no persistent uploads. Preserve aspect ratio; avoid automatically recoloring logos. Maskable icons need safe zones.

Docker Compose binds to loopback to avoid unintentional exposure; Caddy can terminate HTTPS.
