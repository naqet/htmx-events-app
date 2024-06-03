const colors = require('tailwindcss/colors')
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ"],
  theme: {
    extend: {
        colors: {
            "secondary": colors.gray[100] + "4D"
        }
    },
  },
  plugins: [],
}

