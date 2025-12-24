export async function getLocal(key) {
  const v = await chrome.storage.local.get(key);
  return v[key];
}

export async function setLocal(key, value) {
  await chrome.storage.local.set({ [key]: value });
}

export async function removeLocal(key) {
  await chrome.storage.local.remove(key);
}
