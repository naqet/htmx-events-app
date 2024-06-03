const colors = require('tailwindcss/colors')
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ"],
  theme: {
    extend: {
        colors: {
            "secondary": colors.stone[50]
        }
    },
  },
  plugins: [],
}

