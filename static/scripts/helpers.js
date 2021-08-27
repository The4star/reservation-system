const notification = (type, text) => {
  notie.alert({
    type,
    text,
  })
}

const modal = (title, text, icon, confirmButtonText) => {
  Swal.fire({
    title,
    text,
    icon,
    confirmButtonText
  })
}


const toast = (options) => {
  const {
    title = "",
    icon = "success",
    position = "top-end"
  } = options

  const Toast = Swal.mixin({
    toast: true,
    title,
    icon,
    position,
    showConfirmButton: false,
    timer: 3000,
    timerProgressBar: true,
    didOpen: (toast) => {
      toast.addEventListener('mouseenter', Swal.stopTimer)
      toast.addEventListener('mouseleave', Swal.resumeTimer)
    }
  })

  Toast.fire()
}

const datePickerModal = async (options) => {
  const {
    title = "",
    html = ""
  } = options
  const { value: formValues } = await Swal.fire({
    title,
    html,
    focusConfirm: false,
    showCancelButton: true,
    backdrop: false,
    willOpen: () => {
      const rangePicker = document.getElementById('reservation-dates-modal');
      new DateRangePicker(rangePicker, {
        format: "yyyy-mm-dd",
        showOnFocus: true
      });
    },
    preConfirm: () => {
      return [
        document.getElementById('start-date-modal').value,
        document.getElementById('end-date-modal').value
      ]
    },
    didOpen: () => {
      document.getElementById('start-date-modal').removeAttribute('disabled'),
        document.getElementById('end-date-modal').removeAttribute('disabled')
    }
  })

  if (formValues) {
    Swal.fire(JSON.stringify(formValues))
  }
}

export {
  notification,
  modal,
  toast,
  datePickerModal
}