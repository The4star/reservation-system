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

const customModal = (options) => {
  const {
    title = "",
    text = "",
    icon = "",
    showConfirmButton = false,
    showCancelButton = true,
    confirmButtonText = "",
    html = ""
  } = options

  Swal.fire({
    title,
    text,
    icon,
    showConfirmButton,
    confirmButtonText,
    showCancelButton,
    html
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
  const { value: result } = await Swal.fire({
    title,
    html,
    focusConfirm: false,
    showCancelButton: true,
    backdrop: false,
    confirmButtonText: "Check Availability",
    willOpen: () => {
      const rangePicker = document.getElementById('reservation-dates-modal');
      new DateRangePicker(rangePicker, {
        format: "yyyy-mm-dd",
        showOnFocus: true,
        minDate: new Date()
      });
    },
    preConfirm: () => {
      const startDate = document.getElementById('start-date-modal').value
      const endDate = document.getElementById('end-date-modal').value
      if (startDate !== "" && endDate !== "") {
        return {
          startDate,
          endDate
        }
      }
      return false
    },
    didOpen: () => {
      document.getElementById('start-date-modal').removeAttribute('disabled')
      document.getElementById('end-date-modal').removeAttribute('disabled')
    }
  })

  if (result) {
    if (result.startDate !== "") {
      options.callback(result)
    } else {
      options.callback(false)
    }
  }
}

export {
  notification,
  modal,
  toast,
  datePickerModal,
  customModal
}