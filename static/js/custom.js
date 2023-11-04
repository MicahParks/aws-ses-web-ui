'use strict';

const blankUUID = '00000000-0000-0000-0000-000000000000';

// https://stackoverflow.com/a/46181/14797322
const validateEmail = (email) => {
  return String(email)
    .toLowerCase()
    .match(
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
    ) !== null;
};

function modalString(title, subtitle, secondaryButton, primaryButton, icon, color) {
  /*
    bg-red-100   bg-indigo-100   bg-yellow-100   bg-green-100
    text-red-600 text-indigo-600 text-yellow-600 text-green-600
    bg-red-600   bg-indigo-600   bg-yellow-600   bg-green-600
    bg-red-500   bg-indigo-500   bg-yellow-500   bg-green-500
  */
  //language=HTML
  return `
    <div class="relative z-10" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <!--
        Background backdrop, show/hide based on modal state.

        Entering: "ease-out duration-300"
          From: "opacity-0"
          To: "opacity-100"
        Leaving: "ease-in duration-200"
          From: "opacity-100"
          To: "opacity-0"
      -->
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>

      <div class="fixed inset-0 z-10 overflow-y-auto">
        <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
          <!--
            Modal panel, show/hide based on modal state.

            Entering: "ease-out duration-300"
              From: "opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
              To: "opacity-100 translate-y-0 sm:scale-100"
            Leaving: "ease-in duration-200"
              From: "opacity-100 translate-y-0 sm:scale-100"
              To: "opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          -->
          <div
            class="relative transform overflow-hidden rounded-lg bg-white px-4 pb-4 pt-5 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
            <div class="sm:flex sm:items-start">
              <div
                class="mx-auto flex h-16 w-16 flex-shrink-0 items-center justify-center rounded-full bg-${color}-100 sm:mx-0 sm:h-12 sm:w-12">
                <i class="fas ${icon} text-4xl sm:text-2xl text-${color}-600"></i>
              </div>
              <div class="mt-3 text-center sm:ml-4 sm:mt-0 sm:text-left">
                <h3 class="text-base font-semibold leading-6 text-gray-900" id="modal-title">
                  ${title}
                </h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500">
                    ${subtitle}
                  </p>
                </div>
              </div>
            </div>
            <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
              <button type="button"
                      class="inline-flex w-full justify-center rounded-md bg-${color}-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-${color}-500 sm:ml-3 sm:w-auto">
                ${primaryButton}
              </button>
              <button type="button"
                      class="mt-3 inline-flex w-full justify-center rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 sm:mt-0 sm:w-auto">
                ${secondaryButton}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  `.trim();
}

function modal(title, subtitle, icon, color, options = {
  primaryButton: {
    text: '',
    onClick: function (e, m) {
    }
  },
  secondaryButton: {
    text: '',
    onClick: function (e, m) {
    }
  },
  clickOutsideCloses: false,
}) {
  let m = $(modalString(title, subtitle, options.secondaryButton.text, options.primaryButton.text, icon, color));
  let pButton = m.find('button').eq(0);
  let sButton = m.find('button').eq(1);

  // Set on click.
  if (options.primaryButton.text === '') {
    pButton.remove();
  } else {
    pButton.on('click', function (e) {
      options.primaryButton.onClick(e, m);
    });
  }
  sButton.on('click', function (e) {
    options.secondaryButton.onClick(e, m);
  });
  // Set onclick.

  // Set on click outside.
  if (options.clickOutsideCloses) {
    let outside = m.find('.flex.min-h-full.items-end.justify-center.p-4.text-center');
    m.on('click', function (e) {
      let closest = $(e.target).closest('div');
      if (closest.is(outside)) {
        m.remove();
      }
    });
  }
  // Set on click outside.

  m.appendTo('body');
  return m;
}

function dismissModal(e, m) {
  m.remove();
}

function modalDismiss() {
  return {
    clickOutsideCloses: false,
    primaryButton: {
      text: '',
      onClick: dismissModal,
    },
    secondaryButton: {
      text: 'Dismiss',
      onClick: dismissModal,
    },
  }
}

function modalRedirectOrDismiss(buttonText, url) {
  return {
    clickOutsideCloses: false,
    primaryButton: {
      text: buttonText,
      onClick: function (e, m) {
        window.location.href = url;
      }
    },
    secondaryButton: {
      text: 'Dismiss',
      onClick: dismissModal,
    },
  }
}
