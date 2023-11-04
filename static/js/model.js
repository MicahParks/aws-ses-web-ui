// This file was generated with github.com/MicahParks/go-struct-to-js-class
class Compose {
  /**
   * @param {!string} uuid
   * @param {?string} from
   * @param {?string[]} recipientTo
   * @param {?string[]} recipientCC
   * @param {?string[]} recipientBCC
   * @param {!string} subject
   * @param {!string} body
   * @param {!boolean} sesAccepted
   * @param {!Date} created
   */
  constructor(uuid, from, recipientTo, recipientCC, recipientBCC, subject, body, sesAccepted, created) {
    this.uuid = uuid;
    this.from = from;
    this.recipientTo = recipientTo;
    this.recipientCC = recipientCC;
    this.recipientBCC = recipientBCC;
    this.subject = subject;
    this.body = body;
    this.sesAccepted = sesAccepted;
    this.created = created;
  }

  /**
   * @returns {!Compose}
   */
  fromJSON(j) {
    this.uuid = j.uuid;
    this.from = j.from;
    this.recipientTo = repackArray(j.recipientTo, 1, v => v);
    this.recipientCC = repackArray(j.recipientCC, 1, v => v);
    this.recipientBCC = repackArray(j.recipientBCC, 1, v => v);
    this.subject = j.subject;
    this.body = j.body;
    this.sesAccepted = j.sesAccepted;
    this.created = new Date(j.created);
    return this;
  }

  toJSON() {
    return {
      uuid: this.uuid,
      from: this.from,
      recipientTo: this.recipientTo,
      recipientCC: this.recipientCC,
      recipientBCC: this.recipientBCC,
      subject: this.subject,
      body: this.body,
      sesAccepted: this.sesAccepted,
      created: this.created
    };
  }
}

function arrayBufferToBase64(buffer) {
  let binary = '';
  const bytes = new Uint8Array(buffer);
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return window.btoa(binary);
}

function base64ToArrayBuffer(base64) {
  const binaryString = atob(base64);
  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes.buffer;
}

function repackArray(arr, dimensions, fun) {
  if (dimensions === 0) {
    return fun(arr);
  } else {
    return arr.map(element => repackArray(element, dimensions - 1));
  }
}
