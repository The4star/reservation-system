import { notification } from './helpers.js'
const form = document.querySelector('.needs-validation')

const rangePicker = document.getElementById('reservation-dates');
new DateRangePicker(rangePicker, {
  format: "yyyy-mm-dd"
});

form.addEventListener('submit', function (event) {
  if (!form.checkValidity()) {
    event.preventDefault()
    event.stopPropagation()
    notification("error", "Please fill out the fields below")
  }
  form.classList.add('was-validated')
}, false)
