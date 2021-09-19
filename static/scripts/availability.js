import { notification } from './helpers.js'

const form = document.querySelector('.needs-validation')

form.addEventListener('submit', function (event) {
  if (!form.checkValidity()) {
    event.preventDefault()
    event.stopPropagation()
    notification("error", "Please choose your dates below")
  }
  form.classList.add('was-validated')
}, false)


const rangePicker = document.getElementById('reservation-dates');
new DateRangePicker(rangePicker, {
  format: "yyyy-mm-dd",
  minDate: new Date()
});