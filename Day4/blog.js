let data = [];

function postData(event) {
  event.preventDefault();

  let judul = document.getElementById('judul').value;
  let konten = document.getElementById('konten').value;
  let gambar = document.getElementById('gambar-blog').files[0];

  image = URL.createObjectURL(image);

  let blog = {
    judul,
    konten,
    gambar,
    tanggalPosting: new Date(),
    penulis: 'Hydrilla Fragrant',
  };

  data.push(blog);

  renderBlog();
}

function renderBlog() {
  document.getElementById('contents').innerHTML = '';

  for (let index = 0; index < dataBlog.length; index++) {
    document.getElementById('contents').innerHTML += `
    <div class="blog-list-item">
      <div class="blog-image">
        <img src="${dataBlog[index].image}"/>
      </div>
      <div class="blog-content">
        <div class="btn-group">
          <button class="btn-edit">Edit Postingan</button>
          <button class="btn-post">Posting Blog</button>
        </div>
        <h1>
          <a href="blog-detail.html" target="_blank"> ${dataBlog[index].judul} </a>
        </h1>
        <div class="detail-blog-content">13 Oktober 2022 | Hydrilla Fragrant</div>
      </div>
      <p>
      ${dataBlog[index].konten}
      </p>
      </div>
    </div>`;
  }
}
