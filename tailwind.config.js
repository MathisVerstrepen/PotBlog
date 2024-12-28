/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./components/**/*.{html,js,templ}"],
    theme: {
        extend: {
            colors: {
                beige: "var(--c_beige)",
                dark_beige: "var(--c_dark_beige)",
                orange: "var(--c_orange)",
                dark_orange: "var(--c_dark_orange)",
                red: "var(--c_red)",
                green: "var(--c_green)",
                dark_green: "var(--c_dark_green)",
                grey: "var(--c_grey)",
                black: "var(--c_black)",
            },
        },
        fontFamily: {
            sans: ["ClashDisplay-Variable", "sans-serif"],
            serif: ["Archivo-Variable", "serif"],
        },
    },
    plugins: [],
};
