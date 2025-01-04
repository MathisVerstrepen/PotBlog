/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./components/**/*.{html,js,templ}"],
    theme: {
        extend: {
            colors: {
                beige: 'rgb(var(--c_beige)/<alpha-value>)',
                dark_beige: 'rgb(var(--c_dark_beige)/<alpha-value>)',
                orange: 'rgb(var(--c_orange)/<alpha-value>)',
                dark_orange: 'rgb(var(--c_dark_orange)/<alpha-value>)',
                red: 'rgb(var(--c_red)/<alpha-value>)',
                green: 'rgb(var(--c_green)/<alpha-value>)',
                dark_green: 'rgb(var(--c_dark_green)/<alpha-value>)',
                grey: 'rgb(var(--c_grey)/<alpha-value>)',
                black: 'rgb(var(--c_black)/<alpha-value>)',
            },
        },
        fontFamily: {
            sans: ["ClashDisplay-Variable", "sans-serif"],
            serif: ["Archivo-Variable", "serif"],
        },
    },
    plugins: [],
};
