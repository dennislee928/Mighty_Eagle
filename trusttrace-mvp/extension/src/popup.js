import { addEvent, decryptPreviewEvents, buildDailyBundle, exportCommitmentJson, exportAppealZip } from "./lib/ui_actions.js";

const $ = (id) => document.getElementById(id);

async function refresh() {
  const events = await decryptPreviewEvents();
  $("eventsOut").innerHTML = events.length
    ? events.map(e => `<div>• <b>${e.type}</b> ${e.createdAt}: ${escapeHtml(e.preview)}</div>`).join("")
    : "<div>—</div>";
}

function escapeHtml(s) {
  return (s ?? "").replace(/[&<>"']/g, (c) => ({
    "&":"&amp;","<":"&lt;",">":"&gt;","\"":"&quot;","'":"&#39;"
  }[c]));
}

$("addEvent").addEventListener("click", async () => {
  const type = $("eventType").value;
  const text = $("eventText").value.trim();
  if (!text) return;

  await addEvent({ type, text });
  $("eventText").value = "";
  await refresh();
});

$("buildRollup").addEventListener("click", async () => {
  const bundle = await buildDailyBundle();
  $("rootOut").textContent = bundle.merkleRoot;
});

$("exportCommitment").addEventListener("click", async () => {
  const bundle = await buildDailyBundle();
  await exportCommitmentJson(bundle);
});

$("exportAppeal").addEventListener("click", async () => {
  const bundle = await buildDailyBundle();
  await exportAppealZip(bundle);
});

refresh();
