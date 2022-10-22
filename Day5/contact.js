function tampilkan() {
  let nama = document.getElementById('nama').value;
  let email = document.getElementById('email').value;
  let noHp = document.getElementById('noHp').value;
  let perusahaan = document.getElementById('kerja').value;
  let subjek = document.getElementById('subjek').value;
  let pesan = document.getElementById('pesan').value;

  if (nama == '') {
    return alert('isi nama dahulu!');
  } else if (email == '') {
    return alert('isi email dahulu!');
  } else if (noHp == '') {
    return alert('isi nomor dahulu!');
  }

  let penerimaEmail = 'sacrew.jr@gmail.com';
  let a = document.createElement('a');
  a.href = `mailto:${penerimaEmail}?subject:${subjek}&body= Hello, my name is ${nama}, ${pesan}`;
  a.click();
}
