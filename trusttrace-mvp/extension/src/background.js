import { getSettings, setSettings, listEventsForDay, buildDailyBundle } from "./lib/vault.js";

chrome.runtime.onInstalled.addListener(async () => {
  // Create a periodic alarm to remind users to roll up; does not auto-submit transactions.
  chrome.alarms.create("tt_hourly", { periodInMinutes: 60 });
});

chrome.alarms.onAlarm.addListener(async (alarm) => {
  if (alarm.name !== "tt_hourly") return;
  // Placeholder: could notify user if today's bundle not built.
});

chrome.runtime.onMessage.addListener((msg, _sender, sendResponse) => {
  (async () => {
    if (msg?.type === "TT_GET_SETTINGS") {
      sendResponse({ ok: true, data: await getSettings() });
      return;
    }
    if (msg?.type === "TT_SET_SETTINGS") {
      await setSettings(msg.data);
      sendResponse({ ok: true });
      return;
    }
    if (msg?.type === "TT_LIST_EVENTS_TODAY") {
      const today = new Date().toISOString().slice(0, 10);
      const events = await listEventsForDay(today);
      sendResponse({ ok: true, data: events });
      return;
    }
    if (msg?.type === "TT_BUILD_TODAY_BUNDLE") {
      const today = new Date().toISOString().slice(0, 10);
      const bundle = await buildDailyBundle(today);
      sendResponse({ ok: true, data: bundle });
      return;
    }
    sendResponse({ ok: false, error: "unknown message" });
  })();
  return true; // async
});
