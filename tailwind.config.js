/** @type {import('tailwindcss').Config} */
const {fontFamily} = require('tailwindcss/defaultTheme');

module.exports = {
  content: ["./ui/**/*.templ"],
  theme: {
    extend: {
      colors: {
        'steel-dark': '#21212c',
      },
      fontFamily: {
        fira: ['"Fira Mono"', ...fontFamily.sans],
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}

