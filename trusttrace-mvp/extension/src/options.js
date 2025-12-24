import { setPassphrase } from "./lib/vault.js";

const $ = (id) => document.getElementById(id);

$("save").addEventListener("click", async () => {
  const pass = $("passphrase").value;
  if (!pass || pass.length < 10) {
    $("status").textContent = "Passphrase too short (>= 10 chars).";
    return;
  }
  await setPassphrase(pass);
  $("status").textContent = "Saved.";
  $("passphrase").value = "";
});
