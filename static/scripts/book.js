import { notification } from './helpers.js'
const form = document.querySelector('.needs-validation')

form.addEventListener('submit', function (event) {
  if (!form.checkValidity()) {
    event.preventDefault()
    event.stopPropagation()
    notification("error", "Please fill out the fields below")
  }
  form.classList.add('was-validated')
}, false)