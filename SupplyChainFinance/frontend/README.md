# Frontend (Vue3 Engineering)

## Stack
- Vue 3
- Vite
- Vue Router
- Pinia

## Run in development
1. `cd frontend`
2. `npm install`
3. `npm run dev`
4. Open `http://localhost:5173`

Vite proxy forwards `/api`, `/workflow`, `/health` to `http://localhost:8888`.

## Build for backend hosting
1. `cd frontend`
2. `npm install`
3. `npm run build`

Build output is `frontend/dist`.

Backend (`backend/main.go`) will prioritize serving `../frontend/dist`.
If `dist` does not exist, backend falls back to legacy `backend/web` static files.

## Role-based pages
- `business`: enterprise registration, receivable issue/confirm/settle
- `finance`: financing record operations
- `admin`: access both business and finance pages

Role is from login response and passed by header `X-User-Role` for protected APIs.
