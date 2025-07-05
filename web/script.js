async function fetchJSON(url) {
  const res = await fetch(url);
  if (!res.ok) throw new Error('request failed');
  return res.json();
}

async function loadFilters() {
  const manufacturers = await fetchJSON('/api/manufacturers');
  const manufacturerSelect = document.getElementById('manufacturerFilter');
  manufacturerSelect.innerHTML = '<option value="">All Manufacturers</option>' +
    manufacturers.map(m => `<option value="${m.id}">${m.name}</option>`).join('');

  const categories = await fetchJSON('/api/categories');
  const categorySelect = document.getElementById('categoryFilter');
  categorySelect.innerHTML = '<option value="">All Categories</option>' +
    categories.map(c => `<option value="${c.id}">${c.name}</option>`).join('');
}

async function loadModels() {
  const params = new URLSearchParams();
  const search = document.getElementById('search').value;
  const manufacturerId = document.getElementById('manufacturerFilter').value;
  const categoryId = document.getElementById('categoryFilter').value;
  if (search) params.set('search', search);
  if (manufacturerId) params.set('manufacturerId', manufacturerId);
  if (categoryId) params.set('categoryId', categoryId);
  const models = await fetchJSON('/api/models?' + params.toString());
  const tbody = document.querySelector('#carsTable tbody');
  tbody.innerHTML = models.map(m => {
    return `<tr data-id="${m.id}">
      <td><input type="checkbox" class="selectModel" value="${m.id}"></td>
      <td>${m.name}</td><td>${m.year}</td>
      <td>${getManufacturerName(m.manufacturerId)}</td>
      <td>${getCategoryName(m.categoryId)}</td>
      <td><button class="detailBtn" data-id="${m.id}">Details</button></td>
    </tr>`;
  }).join('');
}

function getManufacturerName(id) {
  const opt = document.querySelector(`#manufacturerFilter option[value="${id}"]`);
  return opt ? opt.textContent : '';
}

function getCategoryName(id) {
  const opt = document.querySelector(`#categoryFilter option[value="${id}"]`);
  return opt ? opt.textContent : '';
}

async function showDetails(id) {
  const data = await fetchJSON('/api/models/' + id);
  const div = document.getElementById('details');
  div.innerHTML = `
    <h2>${data.name}</h2>
    <img src="/images/${data.image}" alt="${data.name}">
    <p><strong>Year:</strong> ${data.year}</p>
    <p><strong>Engine:</strong> ${data.specifications.engine}</p>
    <p><strong>Horsepower:</strong> ${data.specifications.horsepower}</p>
    <p><strong>Transmission:</strong> ${data.specifications.transmission}</p>
    <p><strong>Drivetrain:</strong> ${data.specifications.drivetrain}</p>
  `;
  localStorage.setItem('preferredManufacturer', data.manufacturerId);
  loadRecommendations();
}

async function compareSelected() {
  const ids = Array.from(document.querySelectorAll('.selectModel:checked')).map(i => i.value);
  if (ids.length < 2) return;
  const data = await fetchJSON('/api/models/compare?ids=' + ids.join(','));
  const div = document.getElementById('compare');
  let html = '<h2>Comparison</h2><table><thead><tr><th>Name</th><th>Year</th><th>HP</th></tr></thead><tbody>';
  data.forEach(m => {
    html += `<tr><td>${m.name}</td><td>${m.year}</td><td>${m.specifications.horsepower}</td></tr>`;
  });
  html += '</tbody></table>';
  div.innerHTML = html;
}

async function loadRecommendations() {
  const pref = localStorage.getItem('preferredManufacturer');
  if (!pref) return;
  const data = await fetchJSON('/api/recommendations?manufacturerId=' + pref);
  const div = document.getElementById('recommendations');
  div.innerHTML = '<h2>Recommended for You</h2>' +
    '<ul>' + data.map(m => `<li>${m.name}</li>`).join('') + '</ul>';
}

document.getElementById('filterButton').addEventListener('click', loadModels);
document.getElementById('compareBtn').addEventListener('click', compareSelected);

document.querySelector('#carsTable tbody').addEventListener('click', e => {
  if (e.target.classList.contains('detailBtn')) {
    showDetails(e.target.dataset.id);
  }
});

window.addEventListener('DOMContentLoaded', async () => {
  await loadFilters();
  await loadModels();
  loadRecommendations();
});
