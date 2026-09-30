import { useEffect, useState } from 'react'
import { api } from './api'

function ProductCard({ product, onSelect, onAddToCart }) {
  return (
    <div className="glass p-4 flex flex-col gap-3 hover:scale-[1.02] transition-transform">
      <div
        className="h-36 rounded-xl bg-white/20 cursor-pointer"
        onClick={() => onSelect(product.id)}
        aria-label={product.name}
      />
      <div>
        <h3 className="text-white font-semibold">{product.name}</h3>
        <p className="text-white/70 text-sm">${product.price.toFixed(2)}</p>
      </div>
      <button
        className="glass px-3 py-2 text-white text-sm hover:bg-white/20"
        onClick={() => onAddToCart(product.id)}
      >
        Add to cart
      </button>
    </div>
  )
}

function ProductDetail({ productId, onClose }) {
  const [product, setProduct] = useState(null)
  const [reviews, setReviews] = useState([])
  const [comment, setComment] = useState('')
  const [author, setAuthor] = useState('')

  useEffect(() => {
    api.getProduct(productId).then(setProduct).catch(() => {})
    api.listReviews(productId).then(setReviews).catch(() => {})
  }, [productId])

  async function submitReview(e) {
    e.preventDefault()
    if (!author || !comment) return
    const review = await api.addReview(productId, { author, rating: 5, comment })
    setReviews((prev) => [review, ...prev])
    setComment('')
  }

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center p-4 z-20">
      <div className="glass max-w-lg w-full p-6 text-white max-h-[80vh] overflow-y-auto">
        <button className="float-right text-white/70" onClick={onClose}>close</button>
        {product && (
          <>
            <h2 className="text-xl font-bold">{product.name}</h2>
            <p className="text-white/70">${product.price.toFixed(2)}</p>
          </>
        )}
        <h3 className="mt-4 font-semibold">Reviews</h3>
        <ul className="space-y-2 my-2">
          {reviews.map((r) => (
            <li key={r.id} className="text-sm text-white/80 border-b border-white/10 pb-1">
              <strong>{r.author}</strong> ({r.rating}★): {r.comment}
            </li>
          ))}
          {reviews.length === 0 && <li className="text-white/50 text-sm">No reviews yet.</li>}
        </ul>
        <form onSubmit={submitReview} className="flex flex-col gap-2 mt-2">
          <input
            className="glass px-3 py-2 text-sm bg-white/5 placeholder-white/50"
            placeholder="Your name"
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
          />
          <textarea
            className="glass px-3 py-2 text-sm bg-white/5 placeholder-white/50"
            placeholder="Write a review..."
            value={comment}
            onChange={(e) => setComment(e.target.value)}
          />
          <button className="glass px-3 py-2 hover:bg-white/20" type="submit">Post review</button>
        </form>
      </div>
    </div>
  )
}

export default function App() {
  const [products, setProducts] = useState([])
  const [cart, setCart] = useState([])
  const [selectedId, setSelectedId] = useState(null)

  useEffect(() => {
    api.listProducts().then(setProducts).catch(() => {})
    api.getCart().then(setCart).catch(() => {})
  }, [])

  async function addToCart(productId) {
    await api.addToCart(productId)
    setCart(await api.getCart())
  }

  const cartCount = cart.reduce((sum, item) => sum + item.quantity, 0)

  return (
    <div className="p-6 max-w-5xl mx-auto">
      <header className="flex items-center justify-between mb-8">
        <h1 className="text-2xl font-bold text-white">ecom</h1>
        <div className="glass px-4 py-2 text-white">Cart: {cartCount}</div>
      </header>

      <div className="grid grid-cols-2 md:grid-cols-3 gap-5">
        {products.map((p) => (
          <ProductCard key={p.id} product={p} onSelect={setSelectedId} onAddToCart={addToCart} />
        ))}
      </div>

      {selectedId && <ProductDetail productId={selectedId} onClose={() => setSelectedId(null)} />}
    </div>
  )
}
