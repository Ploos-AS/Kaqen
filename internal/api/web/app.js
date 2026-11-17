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
  }
  function choose(file) {
    resetPreview();
    if (!file) { message("No file selected."); return; }
    if (!["image/png", "image/webp"].includes(file.type)) {
      message("Choose a PNG or WebP image. SVG export is not supported yet."); return;
    }
    if (file.size > limit) { message("File exceeds the 5 MiB limit."); return; }
    objectURL = URL.createObjectURL(file);
    const probe = new Image();
    probe.onload = () => {
      if (!probe.width || !probe.height || probe.width * probe.height > 16000000) {
        resetPreview(); message("Image exceeds the 16 million pixel limit."); return;
      }
      selected = file;
      image.src = objectURL;
      filename.textContent = file.name;
      info.textContent = `${probe.width} × ${probe.height} · ${(file.size / 1024).toFixed(1)} KiB`;
      preview.hidden = false;
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
    message("Generating your favicon package…");
    try {
      const data = new FormData();
      data.append("logo", selected, selected.name);
      const response = await fetch("/api/v1/favicons", { method: "POST", body: data });
      if (!response.ok) throw new Error(`Generation failed (HTTP ${response.status}).`);
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
      button.disabled = !selected;
    }
  });
  window.addEventListener("pagehide", () => { if (objectURL) URL.revokeObjectURL(objectURL); });
})();
