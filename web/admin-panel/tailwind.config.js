/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{js,ts,jsx,tsx,mdx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        nebula: {
          50: '#f5f5f7', 100: '#e8e8ed', 200: '#d1d1d6',
          300: '#aeaeb2', 400: '#8e8e93', 500: '#636366',
          600: '#48484a', 700: '#363639', 800: '#2c2c2e',
          900: '#1c1c1e', 950: '#0a0a0a',
        },
        accent: { DEFAULT: '#007AFF', light: '#409CFF', dark: '#0055CC' },
      },
      fontFamily: {
        sans: ['-apple-system', 'BlinkMacSystemFont', 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', 'sans-serif'],
        mono: ['SF Mono', 'SFMono-Regular', 'ui-monospace', 'monospace'],
      },
      borderRadius: { 'apple': '10px', 'apple-lg': '14px' },
    },
  },
  plugins: [],
}
