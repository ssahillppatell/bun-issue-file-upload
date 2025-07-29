import FormData from 'form-data';

const url = 'http://localhost:8080/upload';

const file = Bun.file('./sample.png');
const form = new FormData();
form.append('file', file);
form.append('name', 'Help');

try {
  const response = await fetch(url, {
    method: 'POST',
    body: form,
    headers: form.getHeaders(),
  });
  console.log(response);
} catch (error) {
  console.error(error);
}