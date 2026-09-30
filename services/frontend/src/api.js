const BASE = import.meta.env.VITE_API_URL || '/api'

async function request(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json', ...(options.headers || {}) },
    ...options,
  })
  if (!res.ok) throw new Error(`${options.method || 'GET'} ${path} failed: ${res.status}`)
  if (res.status === 204) return null
  return res.json()
}

export const api = {
  listProducts: () => request('/products'),
  getProduct: (id) => request(`/products/${id}`),
  listReviews: (id) => request(`/products/${id}/reviews`),
  addReview: (id, review) => request(`/products/${id}/reviews`, { method: 'POST', body: JSON.stringify(review) }),
  getRecommendations: (id) => request(`/products/${id}/recommendations`),
  getCart: () => request('/cart'),
  addToCart: (productId, quantity = 1) =>
    request('/cart/items', { method: 'POST', body: JSON.stringify({ product_id: productId, quantity }) }),
  removeFromCart: (productId) => request(`/cart/items/${productId}`, { method: 'DELETE' }),
}
