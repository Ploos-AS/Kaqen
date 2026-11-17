# Public beta acceptance checklist

This checklist requires **manual verification** on the actual HTTPS deployment; CI alone does not certify launch readiness.

## Generator UX
- [ ] Upload valid PNG, WebP and restricted static SVG logos; inspect preview and download ZIP.
- [ ] Verify 16/32/180/192/512 icons, maskable variants, ICO, social card, manifest, HTML and instructions.
- [ ] Try malformed SVG, unsupported formats, corrupt images, over-5-MiB uploads and images exceeding 16 million pixels.
- [ ] Verify 503 busy and 413 oversized-upload messages; retry after load subsides.
- [ ] Test download in Firefox, Chromium and Safari (desktop/mobile where available).

## Accessibility
- [ ] Navigate from the skip link through upload, controls, preview and download using keyboard only.
- [ ] Verify focus is visible and all controls have accessible names.
- [ ] Check live status announcements and busy state with a screen reader.
- [ ] Test 200% zoom, narrow mobile viewport, high contrast and reduced motion.
- [ ] Run an automated accessibility audit and manually review findings.

## Privacy and operations
- [ ] Replace placeholder operator contact with a verified contact method.
- [ ] Review the privacy notice against actual server, reverse proxy and CDN logging, including IP retention.
- [ ] Verify multipart temporary-file cleanup and no persistent storage of uploaded logos.
- [ ] Confirm HTTPS, reverse proxy body/time limits, CSP and response cache behavior.
- [ ] Load-test concurrent conversions and measure memory/CPU on the chosen VPS.
- [ ] Check /healthz and /readyz from the monitoring environment.
- [ ] Review backups, deployment rollback and incident response before public exposure.
- [ ] Integrate with Ploos-Deploy only after that project is ready.

Record test date, browser/device, deployment commit and any unresolved issues before opening public beta.
