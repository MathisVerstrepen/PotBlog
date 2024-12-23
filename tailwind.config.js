/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./components/**/*.{html,js,templ}"],
    theme: {
        extend: {
            colors: {
                beige: "#ECDFCC",
                orange: "#DB7E21",
                green: "#697565",
                grey: "#3C3D37",
                black: "#181C14",
            },
        },
    },
    plugins: [],
};
