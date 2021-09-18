import { notification, modal, customModal } from './helpers.js'
import { datePickerModal } from './helpers.js'
const checkAvailabilityButton = document.querySelector('#check-availability')
const html = `
    <form action="" method="GET" class="needs-validation" novalidate>
          <div class="row" id="reservation-dates-modal">
            <div class="col">
              <div class="mb-3">
                <input disabled required type="text" class="form-control" id="start-date-modal" name="start-date" autocomplete="off"
                  placeholder="Select your start date">
              </div>
            </div>
            <div class="col">
              <div class="mb-3">
                <input disabled required type="text" class="form-control" id="end-date-modal" name="end-date" autocomplete="off"
                  placeholder="Select your end date">
              </div>
            </div>
          </div>
        </form>
    `


const openModal = () => {
  datePickerModal({
    html,
    title: "Choose your dates",
    callback: async (result) => {
      const roomType = document.querySelector('#room-type')
      const roomID = roomType.dataset.roomid
      const roomName = roomType.innerHTML
      result.roomID = roomID
      try {
        const response = await axios.post("/room-availability", result)
        const data = response.data;
        if (data.ok) {
          customModal({
            title: `The ${roomName} is available`,
            text: `Would you like to book the ${roomName}`,
            icon: "success",
            showConfirmationButton: false,
            showCancelButton: true,
            html: `
              <p>
                <a href="/book-room?id=${roomID}&sd=${result.startDate}&ed=${result.endDate}" class="btn btn-primary mt-1">Book Now</a>
              </p>
            `
          })
        } else {
          modal("Room is not available", "Please try alternative dates", "error", "OK")
        }
      } catch (error) {
        console.log(error);
      }
    }
  })
  setTimeout(() => {
    const form = document.querySelector('.needs-validation')
    const confirmButton = document.querySelector('.swal2-confirm')
    confirmButton.addEventListener('click', function (event) {
      if (!form.checkValidity()) {
        event.preventDefault()
        event.stopPropagation()
        notification("error", "Please fill out the fields below")
      }
      form.classList.add('was-validated')
    }, false)
  }, 200);
}

checkAvailabilityButton.addEventListener('click', openModal)