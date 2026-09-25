export async function request(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) }
  const user = JSON.parse(localStorage.getItem('scf_user') || 'null')
  const role = String(user?.role || '').trim().toLowerCase()
  if (options.withRole && role) {
    headers['X-User-Role'] = role
  }

  const resp = await fetch(path, {
    method: options.method || 'GET',
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined
  })

  const data = await resp.json().catch(() => ({}))
  if (!resp.ok) {
    throw data
  }
  return data
}

export async function uploadFiles(path, files, withRole = false, fields = {}) {
  const form = new FormData()
  for (const file of files) {
    form.append('files', file)
  }
  Object.entries(fields || {}).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      form.append(key, String(value))
    }
  })

  const headers = {}
  const user = JSON.parse(localStorage.getItem('scf_user') || 'null')
  const role = String(user?.role || '').trim().toLowerCase()
  if (withRole && role) {
    headers['X-User-Role'] = role
  }

  const resp = await fetch(path, {
    method: 'POST',
    headers,
    body: form
  })
  const data = await resp.json().catch(() => ({}))
  if (!resp.ok) {
    throw data
  }
  return data
}

export async function requestPreferDB(dbPath, fallbackPath, options = {}) {
  try {
    return await request(dbPath, options)
  } catch (_) {
    return request(fallbackPath, options)
  }
}
