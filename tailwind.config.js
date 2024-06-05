const colors = require('tailwindcss/colors')
/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./views/**/*.templ"],
    theme: {
        extend: {
            colors: {
                "secondary": colors.stone[50]
            },
            keyframes: {
                appear: {
                    "0%": {
                        opacity: 0,
                        transform: `translateX(10px) scale(0.95)`,
                    },

                    "100%": {
                        opacity: 100,
                        transform: `translateX(0) scale(1)`,
                    }
                }
            },
            animation: {
                appear: 'appear 300ms'
            }
        },
    },
    plugins: [],
    safelist: [
        "animation-appear",
        "toast",
        "toast-success",
        "toast-info",
        "toast-danger",
        "toast-message",
    ]
}

