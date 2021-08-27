import { datePickerModal } from './helpers.js'
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

const checkAvailabilityButton = document.querySelector('#check-availability')
checkAvailabilityButton.addEventListener('click', () => datePickerModal({ html, title: "Choose your dates" }))