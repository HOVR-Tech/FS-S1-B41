const contactForm = document.getElementById('contact-form');

contactForm.addEventListener('submit', () => {
  sendMail();
});

const sendMail = () => {
  let name = document.getElementById('name').value;
  let email = document.getElementById('email').value;
  let phone = document.getElementById('phone').value;
  let subject = document.getElementById('subject').value;
  let message = document.getElementById('message').value;

  if (name === '') {
    return alert('Name is required');
  }
  if (email === '') {
    return alert('email is required');
  }
  if (phone === '') {
    return alert('Phone is required');
  }
  if (subject === '') {
    return alert('Subject is required');
  }

  const emailReciever = 'hydrilla.salim@gmail.com';

  const a = document.createElement('a');

  a.href = `mailto:${emailReciever}?subject=${subject}&body= Hello, my name is ${name}, please contact me at ${phone}, ${message}`;
  console.log(a);
  a.click();

  alert('Your Message was Sent Successfully, Please Wait for My Reply');
};
