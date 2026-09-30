/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,jsx}'],
  theme: {
    extend: {
      backgroundImage: {
        'app-gradient': 'linear-gradient(135deg, #4f46e5 0%, #7c3aed 45%, #db2777 100%)',
      },
    },
  },
  plugins: [],
}
