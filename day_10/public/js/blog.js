let dataBlog = [];

function addBlog(event) {
  event.preventDefault();

  let title = document.getElementById('input-title').value;
  let content = document.getElementById('input-content').value;
  let image = document.getElementById('input-blog-image').files[0];

  image = URL.createObjectURL(image);

  let blog = {
    title,
    content,
    image,
    postAt: new Date(),
    author: 'Hydrilla Fragrant',
  };

  dataBlog.push(blog);
  console.log(dataBlog);

  renderBlog();
}

function renderBlog() {
  document.getElementById('contents').innerHTML = '';

  for (let index = 0; index < dataBlog.length; index++) {
    document.getElementById('contents').innerHTML += `
    <div class="blog-list-item">
    <div class="blog-image">
        <img src="${dataBlog[index].image}">
    </div>
    <div class="blog-content">
        <div class="btn-group">
            <button class="btn-edit">Edit Post</button>
            <button class="btn-post">Post Blog</button>
        </div>
        <h1>
            <a href="blog-detail.html" target="_blank">
                ${dataBlog[index].title}
            </a>
        </h1>
        <div class="detail-blog-content">
            ${getFullTime(dataBlog[index].postAt)} | ${dataBlog[index].author}
        </div>
        <p>
            ${dataBlog[index].content}
        </p>
        <div>
            <p style="font-size: 15px; color: grey">${getDistanceTime(dataBlog[index].postAt)}</p>
        </div>
    </div>
</div>`;
  }
}
function getFullTime(time) {
  // time = new Date();

  let monthName = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  let date = time.getDate();

  let monthIndex = time.getMonth();
  let year = time.getFullYear();

  let hours = time.getHours();
  let minutes = time.getMinutes();

  if (hours <= 9) {
    hours = '0' + hours;
  } else if (minutes <= 9) {
    minutes = '0' + minutes;
  }

  return `${date} ${monthName[monthIndex]} ${year} ${hours}:${minutes} WIB`;
}

function getDistanceTime(time) {
  let timeNow = new Date();
  let timePost = time;

  let distance = timeNow - timePost;

  let miliSecond = 1000;
  let secondInHours = 3600;
  let hoursInDay = 24;

  let distanceDay = Math.floor(distance / (miliSecond * secondInHours * hoursInDay));
  let distanceHours = Math.floor(distance / (miliSecond * 60 * 60));
  let distanceMinutes = Math.floor(distance / (miliSecond * 60));
  let distanceSecond = Math.floor(distance / miliSecond);

  if (distanceDay > 0) {
    return `${distanceDay} day(s) ago`;
  } else if (distanceHours > 0) {
    return `${distanceHours} hour(s) ago`;
  } else if (distanceMinutes > 0) {
    return `${distanceMinutes} minute(s) ago`;
  } else {
    return `${distanceSecond} second(s) ago`;
  }
}

setInterval(function () {
  renderBlog();
}, 5000);

// setInterval(interval, 3000);

// function interval() {
//   renderBlog();
// }
