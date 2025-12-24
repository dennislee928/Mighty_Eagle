# Extension (MV3)

## Notes
- This MVP stores encrypted events in `chrome.storage.local` for simplicity.
- Upgrade path: move encrypted records to IndexedDB; keep settings in storage.

## Build & load
```bash
npm install
npm run build
```

Chrome → `chrome://extensions` → Developer mode → Load unpacked → select `extension/dist`.

## Guardrails
- No scraping and no automated actions on Tinder.
- User-entered notes only.
