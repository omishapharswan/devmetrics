# DevMetrics frontend

Next.js (App Router) dashboard for DevMetrics. See the [repository root
README](../README.md) for the full setup guide, architecture, and API
reference — this file only covers commands specific to this package.

```bash
npm install
npm run dev     # dev server on http://localhost:3000
npm run build   # production build
npm run lint    # eslint
npm run test    # vitest (unit + component tests)
```

The dashboard expects the backend API at `http://localhost:8080` by
default; override with `NEXT_PUBLIC_API_BASE_URL` (see `.env.example`).
