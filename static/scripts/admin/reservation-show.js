import { toast, customModal } from '../helpers.js';
const form = document.querySelector('.needs-validation');
const processButton = document.querySelector("#process");
const deleteButton = document.querySelector("#delete")
const processRes = () => {
  const { id, src, year, month } = processButton.dataset

  customModal({
    title: `Are you sure?`,
    icon: "warning",
    showConfirmationButton: false,
    showCancelButton: true,
    html: `
      <p>
        <a href="/admin/reservations/${src}/${id}/process?y=${year}&m=${month}" class="btn btn-primary mt-1">OK</a>
      </p>
    `
  })
}

const deleteRes = () => {
  const { id, src, year, month } = deleteButton.dataset
  customModal({
    title: `Are you sure?`,
    icon: "warning",
    showConfirmationButton: false,
    showCancelButton: true,
    html: `
      <p>
        <a href="/admin/reservations/${src}/${id}/delete?y=${year}&m=${month}" class="btn btn-primary mt-1">OK</a>
      </p>
    `
  })
}

form.addEventListener('submit', function (event) {
  if (!form.checkValidity()) {
    event.preventDefault()
    event.stopPropagation()
    toast({
      icon: "error",
      title: "Please address the errors in the form",
      position: "bottom-end"
    })
  }
  form.classList.add('was-validated')
}, false)

if (processButton) {
  processButton.addEventListener("click", processRes);
}
deleteButton.addEventListener("click", deleteRes)