/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../templates/*.{gohtml,html,js}","../sripts/*.js"],
  theme: {
    container:{
      center:true,
    },
    extend: {},
  },
  plugins: [require("daisyui"),require("tailwind"),require("autoprefixer")],
}
