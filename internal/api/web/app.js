(() => {
  "use strict";
  const form = document.querySelector("#generator");
  const input = document.querySelector("#logo");
  const drop = document.querySelector("#drop");
  const preview = document.querySelector("#preview");
  const image = document.querySelector("#preview-image");
  const filename = document.querySelector("#file-name");
  const info = document.querySelector("#file-info");
  const status = document.querySelector("#status");
  const button = document.querySelector("#generate");
  const padding = document.querySelector("#padding");
  const paddingValue = document.querySelector("#padding-value");
  const background = document.querySelector("#background");
  const customColor = document.querySelector("#custom-color");
  const customLabel = document.querySelector("#custom-color-label");
  const outputPreview = document.querySelector("#output-preview");
  const canvases = document.querySelectorAll("[data-icon-size]");
  let decodedImage = null;

  function selectedBackground() {
    return background.value === "custom" ? customColor.value : background.value;
  }
  function renderPreview() {
    if (!decodedImage) return;
    const fraction = Number(padding.value) / 100;
    const bg = selectedBackground();
    for (const canvas of canvases) {
      const ctx = canvas.getContext("2d");
      const size = canvas.width;
      ctx.clearRect(0, 0, size, size);
      if (bg !== "transparent") {
        ctx.fillStyle = bg;
        ctx.fillRect(0, 0, size, size);
      }
      const available = Math.max(1, Math.floor(size * (1 - 2 * fraction)));
      const width = decodedImage.naturalWidth;
      const height = decodedImage.naturalHeight;
      const targetWidth = width > height ? available : Math.max(1, Math.floor(available * width / height));
      const targetHeight = width > height ? Math.max(1, Math.floor(available * height / width)) : available;
      ctx.imageSmoothingEnabled = true;
      ctx.imageSmoothingQuality = "high";
      ctx.drawImage(decodedImage, Math.floor((size - targetWidth) / 2), Math.floor((size - targetHeight) / 2), targetWidth, targetHeight);
    }
    renderMaskablePreviews();
    outputPreview.hidden = false;
  }
  function renderMaskablePreviews() {
    const size = 192;
    const bg = selectedBackground() === "transparent" ? "#ffffff" : selectedBackground();
    const width = decodedImage.naturalWidth;
    const height = decodedImage.naturalHeight;
    const available = Math.floor(size * 0.54);
    const targetWidth = width > height ? available : Math.max(1, Math.floor(available * width / height));
    const targetHeight = width > height ? Math.max(1, Math.floor(available * height / width)) : available;
    for (const [id, rounded] of [["mask-circle", false], ["mask-rounded", true]]) {
      const canvas = document.querySelector("#" + id);
      const ctx = canvas.getContext("2d");
      ctx.clearRect(0, 0, size, size);
      ctx.save();
      ctx.beginPath();
      if (rounded) {
        ctx.roundRect(0, 0, size, size, size * 0.22);
      } else {
        ctx.arc(size / 2, size / 2, size / 2, 0, Math.PI * 2);
      }
      ctx.clip();
      ctx.fillStyle = bg;
      ctx.fillRect(0, 0, size, size);
      ctx.imageSmoothingEnabled = true;
      ctx.imageSmoothingQuality = "high";
      ctx.drawImage(decodedImage, Math.floor((size - targetWidth) / 2), Math.floor((size - targetHeight) / 2), targetWidth, targetHeight);
      ctx.restore();
    }
  }
  function updateColorControl() {
    const custom = background.value === "custom";
    customColor.hidden = !custom;
    customLabel.hidden = !custom;
    renderPreview();
  }
  padding.addEventListener("input", () => {
    paddingValue.textContent = padding.value + "%";
    renderPreview();
  });
  background.addEventListener("change", updateColorControl);
  customColor.addEventListener("input", renderPreview);
  let selected = null;
  let objectURL = null;
  const limit = 5 * 1024 * 1024;

  function message(text) { status.textContent = text; }
  function resetPreview() {
    if (objectURL) URL.revokeObjectURL(objectURL);
    objectURL = null;
    selected = null;
    button.disabled = true;
    preview.hidden = true;
    outputPreview.hidden = true;
    decodedImage = null;
  }
  function choose(file) {
    resetPreview();
    if (!file) { message("No file selected."); return; }
    if (!["image/png", "image/webp", "image/svg+xml"].includes(file.type)) {
      message("Choose a PNG, WebP or SVG image."); return;
    }
    if (file.size > limit) { message("File exceeds the 5 MiB limit."); return; }
    objectURL = URL.createObjectURL(file);
    const probe = new Image();
    probe.onload = () => {
      if (!probe.width || !probe.height || probe.width * probe.height > 16000000) {
        resetPreview(); message("Image exceeds the 16 million pixel limit."); return;
      }
      selected = file;
      decodedImage = probe;
      image.src = objectURL;
      filename.textContent = file.name;
      info.textContent = `${probe.width} × ${probe.height} · ${(file.size / 1024).toFixed(1)} KiB`;
      preview.hidden = false;
      renderPreview();
      button.disabled = false;
      message("Ready to generate.");
    };
    probe.onerror = () => { resetPreview(); message("Cannot read this image."); };
    probe.src = objectURL;
  }
  input.addEventListener("change", () => choose(input.files[0]));
  for (const eventName of ["dragenter", "dragover"]) {
    drop.addEventListener(eventName, event => {
      event.preventDefault();
      drop.classList.add("dragging");
    });
  }
  drop.addEventListener("dragleave", () => drop.classList.remove("dragging"));
  drop.addEventListener("drop", event => {
    event.preventDefault();
    drop.classList.remove("dragging");
    choose(event.dataTransfer.files[0]);
  });
  form.addEventListener("submit", async event => {
    event.preventDefault();
    if (!selected) return;
    button.disabled = true;
    form.setAttribute("aria-busy", "true");
    message("Generating your favicon package…");
    try {
      const data = new FormData();
      data.append("logo", selected, selected.name);
      data.append("padding", padding.value);
      data.append("background", selectedBackground());
      const response = await fetch("/api/v1/favicons", { method: "POST", body: data });
      if (!response.ok) {
        if (response.status === 503) throw new Error("Server is busy. Please try again shortly.");
        if (response.status === 413) throw new Error("Upload exceeds the size limit.");
        throw new Error(`Generation failed (HTTP ${response.status}).`);
      }
      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      const anchor = document.createElement("a");
      anchor.href = url;
      anchor.download = "kaqen-favicons.zip";
      document.body.appendChild(anchor);
      anchor.click();
      anchor.remove();
      setTimeout(() => URL.revokeObjectURL(url), 60000);
      message("Your ZIP is ready. Download should begin automatically.");
    } catch (error) {
      message(error.message || "Unable to generate your ZIP.");
    } finally {
      form.setAttribute("aria-busy", "false");
      button.disabled = !selected;
    }
  });
  window.addEventListener("pagehide", () => { if (objectURL) URL.revokeObjectURL(objectURL); });
})();
