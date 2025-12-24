/**
 * Content script intentionally minimal:
 * - Only sets a flag that the current tab is tinder.com for UX context.
 * - Does NOT extract profile/chat data.
 */
chrome.runtime.sendMessage({ type: "TT_PAGE_CONTEXT", data: { host: location.host } }).catch(() => {});
