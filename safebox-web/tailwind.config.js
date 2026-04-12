/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#3B82F6',
          dark: '#2563EB',
          darker: '#1D4ED8',
        },
        safe: {
          bg: '#F0F4F8',
          card: '#FFFFFF',
          dark: '#1E293B',
        },
      },
      fontFamily: {
        sans: ['Inter', 'PingFang SC', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
