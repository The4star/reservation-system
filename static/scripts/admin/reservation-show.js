import { notification } from '../helpers.js'
const form = document.querySelector('.needs-validation')
console.log(form);
form.addEventListener('submit', function (event) {
  if (!form.checkValidity()) {
    event.preventDefault()
    event.stopPropagation()
    toast({
      icon: "error",
      title: "Please address the errors in the form",
      position: "center-end"
    })
  }
  form.classList.add('was-validated')
}, false)