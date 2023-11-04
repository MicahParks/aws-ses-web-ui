/** @type {import('tailwindcss').Config} */
const defaultTheme = require('tailwindcss/defaultTheme');
module.exports = {
  content: ["./templates/**/*.gohtml", "./static/js/**/*.js"],
  theme: {
    extend: {},
    fontFamily: {
      'sans': [...defaultTheme.fontFamily.sans, '"Font Awesome 6 Pro"']
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}