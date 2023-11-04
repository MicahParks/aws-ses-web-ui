$(function () {
  let rTo = [];
  let rCc = [];
  let rBcc = [];
  let hasSubject = false;
  let hasBody = false;

  let from = $('#from');
  let recipientsTo = $('#recipient-to');
  let recipientsCc = $('#recipient-cc');
  let recipientsBcc = $('#recipient-bcc');
  let sendIcon = $('#send-icon');

  function complete(jqXHR, status) {
    switch (jqXHR.status) {
      case 202:
        modal(
          'Email sent',
          `The email has been sent to <strong>${rTo.length + rCc.length + rBcc.length}</strong> recipients.`,
          'fa-envelope-circle-check', 'green',
          modalDismiss(),
        );
        subject.val('');
        body.val('');
        recipientsTo.val('');
        recipientsCc.val('');
        recipientsBcc.val('');
        rTo.splice(0, rTo.length);
        rCc.splice(0, rCc.length);
        rBcc.splice(0, rBcc.length);
        updateSendButton();
        break;
      default:
        modal(
          'Unknown error',
          `An unknown error occurred.`,
          'fa-question-circle', 'red',
          modalDismiss(),
        );
        break;
    }
    sendIcon.removeClass('fa-spinner fa-spin').addClass('fa-paper-plane');
  }

  let subject = $('#subject');
  let body = $('#body');
  let sendButton = $('#send-button');
  sendButton.on('click', function (e) {
    e.preventDefault();

    sendButton.attr('disabled', '');

    modal(
      'Confirmation',
      `Are you sure you want to send this email to <strong>${rTo.length + rCc.length + rBcc.length}</strong> recipients?`,
      'fa-thumbs-up', 'green',
      {
        clickOutsideCloses: true,
        secondaryButton: {
          onClick: dismissModal,
          text: 'Cancel',
        },
        primaryButton: {
          onClick: function (e, m) {
            let c = new Compose(blankUUID, from.val(), rTo, rCc, rBcc, subject.val(), body.val(), false, new Date());
            $.ajax('/api/compose', {
              accepts: 'application/json',
              contentType: 'application/json',
              data: JSON.stringify(c),
              dataType: 'json',
              method: 'POST',
              complete: complete,
            });
            sendIcon.removeClass('fa-paper-plane').addClass('fa-spinner fa-spin');
            dismissModal(e, m);
          },
          text: 'Send'
        }
      }
    );

    return false;
  });

  function updateSendButton() {
    if (from.val() !== "" && ((rTo.length > 0) || (rCc.length > 0) || (rBcc.length > 0)) && hasSubject && hasBody) {
      sendButton.removeAttr('disabled');
    } else {
      sendButton.attr('disabled', '');
    }
  }

  function updateRecipientArray(array, input) {
    let v = input.val().trim();
    if (v.length > 0) {
      let emails = v.split(',');
      for (let i = 0; i < emails.length; i++) {
        let email = emails[i].trim();
        if (email.length > 0) {
          array.push(email);
        }
      }
    }
  }

  from.on('input', function () {
    updateSendButton();
  });
  recipientsTo.on('input', function () {
    rTo.splice(0, rTo.length);
    updateRecipientArray(rTo, recipientsTo);
    updateSendButton();
  });
  recipientsCc.on('input', function () {
    rCc.splice(0, rCc.length);
    updateRecipientArray(rCc, recipientsCc);
    updateSendButton();
  });
  recipientsBcc.on('input', function () {
    rBcc.splice(0, rBcc.length);
    updateRecipientArray(rBcc, recipientsBcc);
    updateSendButton();
  });
  subject.on('input', function () {
    hasSubject = subject.val().trim().length > 0;
    updateSendButton();
  });
  body.on('input', function () {
    body.css('height', body[0].scrollHeight + 'px');
    hasBody = body.val().trim().length > 0;
    updateSendButton();
  });
});