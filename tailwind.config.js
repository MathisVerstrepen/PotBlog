/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./components/**/*.{html,js,templ}"],
    theme: {
        extend: {
            colors: {
                beige: "#ECDFCC",
                beige_dark: "#D1BFA8",
                orange: "#DB7E21",
                green: "#697565",
                dark_green: "#4e564b",
                grey: "#3C3D37",
                black: "#181C14",
            },
        },
        fontFamily: {
            sans: ["ClashDisplay-Variable", "sans-serif"],
            serif: ["Archivo-Variable", "serif"],
        },
    },
    plugins: [],
};
